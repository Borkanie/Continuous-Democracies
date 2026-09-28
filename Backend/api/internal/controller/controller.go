package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/borkanie/brand-new-day-api/internal/generated"
	"github.com/borkanie/brand-new-day-api/internal/repository"
	"github.com/borkanie/brand-new-day-api/internal/service"
)

// Controller implements the generated.ServerInterface, delegating every request to the
// appropriate service. It performs parameter parsing and response writing only -- no business
// logic lives here.
type Controller struct {
	politicianService *service.PoliticianService
	votingService     *service.VotingService
	lawService        *service.LawService
}

// NewController constructs a Controller from its service dependencies.
func NewController(
	politicianService *service.PoliticianService,
	votingService *service.VotingService,
	lawService *service.LawService,
) *Controller {
	return &Controller{
		politicianService: politicianService,
		votingService:     votingService,
		lawService:        lawService,
	}
}

// Compile-time assertion that Controller satisfies the generated ServerInterface.
var _ generated.ServerInterface = (*Controller)(nil)

// writeJSONResponse writes payload as a JSON body with the given status code, substituting an
// empty slice/object as needed so list endpoints never serialize a nil slice as "null".
func writeJSONResponse(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	if err := json.NewEncoder(responseWriter).Encode(payload); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

// writeErrorResponse writes a generated.ErrorResponse body with the given status code.
func writeErrorResponse(responseWriter http.ResponseWriter, statusCode int, message string) {
	writeJSONResponse(responseWriter, statusCode, generated.ErrorResponse{Error: message})
}

// handleServiceError inspects a service-layer error and writes the appropriate HTTP response:
// repository.ErrNotFound becomes a 404, anything else is logged and becomes a 500.
func handleServiceError(responseWriter http.ResponseWriter, request *http.Request, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		writeErrorResponse(responseWriter, http.StatusNotFound, "not found")
		return
	}
	slog.ErrorContext(request.Context(), "unhandled service error", "path", request.URL.Path, "error", err)
	writeErrorResponse(responseWriter, http.StatusInternalServerError, "internal server error")
}

// GetHealth returns a static ok status and touches no service.
func (controller *Controller) GetHealth(responseWriter http.ResponseWriter, request *http.Request) {
	writeJSONResponse(responseWriter, http.StatusOK, generated.HealthResponse{Status: "ok"})
}

// ListParties handles GET /parties.
func (controller *Controller) ListParties(responseWriter http.ResponseWriter, request *http.Request) {
	partyDTOs, err := controller.politicianService.ListParties(request.Context())
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	if partyDTOs == nil {
		partyDTOs = []generated.Party{}
	}
	writeJSONResponse(responseWriter, http.StatusOK, partyDTOs)
}

// GetPartyById handles GET /parties/{partyId}.
func (controller *Controller) GetPartyById(responseWriter http.ResponseWriter, request *http.Request, partyId string) {
	partyDTO, err := controller.politicianService.GetPartyByID(request.Context(), partyId)
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	writeJSONResponse(responseWriter, http.StatusOK, partyDTO)
}

// ListPoliticians handles GET /politicians.
func (controller *Controller) ListPoliticians(responseWriter http.ResponseWriter, request *http.Request) {
	politicianDTOs, err := controller.politicianService.ListPoliticians(request.Context())
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	if politicianDTOs == nil {
		politicianDTOs = []generated.Politician{}
	}
	writeJSONResponse(responseWriter, http.StatusOK, politicianDTOs)
}

// GetPoliticianById handles GET /politicians/{politicianId}.
func (controller *Controller) GetPoliticianById(responseWriter http.ResponseWriter, request *http.Request, politicianId string) {
	politicianDTO, err := controller.politicianService.GetPoliticianByID(request.Context(), politicianId)
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	writeJSONResponse(responseWriter, http.StatusOK, politicianDTO)
}

// GetVotesByPolitician handles GET /politicians/{politicianId}/votes.
func (controller *Controller) GetVotesByPolitician(responseWriter http.ResponseWriter, request *http.Request, politicianId string) {
	voteEntryDTOs, err := controller.politicianService.GetVotesByPolitician(request.Context(), politicianId)
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	if voteEntryDTOs == nil {
		voteEntryDTOs = []generated.PoliticianVoteEntry{}
	}
	writeJSONResponse(responseWriter, http.StatusOK, voteEntryDTOs)
}

// ListVotingRounds handles GET /votingRounds.
func (controller *Controller) ListVotingRounds(responseWriter http.ResponseWriter, request *http.Request) {
	votingRoundDTOs, err := controller.votingService.ListVotingRounds(request.Context())
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	if votingRoundDTOs == nil {
		votingRoundDTOs = []generated.VotingRound{}
	}
	writeJSONResponse(responseWriter, http.StatusOK, votingRoundDTOs)
}

// GetVotingRoundById handles GET /votingRounds/{roundId}.
func (controller *Controller) GetVotingRoundById(responseWriter http.ResponseWriter, request *http.Request, roundId int) {
	votingRoundDetailDTO, err := controller.votingService.GetVotingRoundByID(request.Context(), roundId)
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	writeJSONResponse(responseWriter, http.StatusOK, votingRoundDetailDTO)
}

// GetVotesByVotingRound handles GET /votingRounds/{roundId}/votes.
func (controller *Controller) GetVotesByVotingRound(responseWriter http.ResponseWriter, request *http.Request, roundId int) {
	hydratedVoteDTOs, err := controller.votingService.GetVotesByVotingRound(request.Context(), roundId)
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	if hydratedVoteDTOs == nil {
		hydratedVoteDTOs = []generated.HydratedVote{}
	}
	writeJSONResponse(responseWriter, http.StatusOK, hydratedVoteDTOs)
}

// ListLawBuckets handles GET /lawBuckets.
func (controller *Controller) ListLawBuckets(responseWriter http.ResponseWriter, request *http.Request) {
	lawBucketDTOs, err := controller.lawService.ListLawBuckets(request.Context())
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	if lawBucketDTOs == nil {
		lawBucketDTOs = []generated.LawBucket{}
	}
	writeJSONResponse(responseWriter, http.StatusOK, lawBucketDTOs)
}

// GetLawBucketById handles GET /lawBuckets/{lawBucketId}.
func (controller *Controller) GetLawBucketById(responseWriter http.ResponseWriter, request *http.Request, lawBucketId int) {
	lawBucketDTO, err := controller.lawService.GetLawBucketByID(request.Context(), lawBucketId)
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	writeJSONResponse(responseWriter, http.StatusOK, lawBucketDTO)
}

// GetVotingRoundsByLawBucket handles GET /lawBuckets/{lawBucketId}/votingRounds.
func (controller *Controller) GetVotingRoundsByLawBucket(responseWriter http.ResponseWriter, request *http.Request, lawBucketId int) {
	votingRoundDTOs, err := controller.lawService.GetVotingRoundsByLawBucket(request.Context(), lawBucketId)
	if err != nil {
		handleServiceError(responseWriter, request, err)
		return
	}
	if votingRoundDTOs == nil {
		votingRoundDTOs = []generated.VotingRound{}
	}
	writeJSONResponse(responseWriter, http.StatusOK, votingRoundDTOs)
}
