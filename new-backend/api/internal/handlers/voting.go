package handlers

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/borkanie/parliament-api/internal/repository"
)

type VotingHandler struct {
	repo *repository.VotingRepo
}

func NewVotingHandler(repo *repository.VotingRepo) *VotingHandler {
	return &VotingHandler{repo: repo}
}

// GET /api/Voting/getAllRounds?startDate?&endDate?&keywords[]?&maxNumberOfEntries=100
func (handler *VotingHandler) GetAllRounds(responseWriter http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()

	roundFilter := repository.RoundFilter{Limit: 100}

	if maxEntriesStr := queryParams.Get("maxNumberOfEntries"); maxEntriesStr != "" {
		count, err := strconv.Atoi(maxEntriesStr)
		if err != nil || count <= 0 {
			writeBadRequest(responseWriter, "invalid 'maxNumberOfEntries' parameter")
			return
		}
		roundFilter.Limit = count
	}

	if startDateStr := queryParams.Get("startDate"); startDateStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			writeBadRequest(responseWriter, "invalid 'startDate': use RFC3339 format")
			return
		}
		roundFilter.StartDate = &parsedTime
	}

	if endDateStr := queryParams.Get("endDate"); endDateStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			writeBadRequest(responseWriter, "invalid 'endDate': use RFC3339 format")
			return
		}
		roundFilter.EndDate = &parsedTime
	}

	// Accept both ?keywords=foo and ?keywords[]=foo (frontend sends single encoded value)
	roundFilter.Keywords = collectKeywords(queryParams)

	rounds, err := handler.repo.GetAllRounds(request.Context(), roundFilter)
	if err != nil {
		writeInternalError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, rounds)
}

// GET /api/Voting/getRoundById/?voteNumber={int}
func (handler *VotingHandler) GetRoundById(responseWriter http.ResponseWriter, request *http.Request) {
	voteNumberStr := request.URL.Query().Get("voteNumber")
	if voteNumberStr == "" {
		writeBadRequest(responseWriter, "missing 'voteNumber' parameter")
		return
	}

	voteId, err := strconv.Atoi(voteNumberStr)
	if err != nil {
		writeBadRequest(responseWriter, "invalid 'voteNumber': must be an integer")
		return
	}

	round, err := handler.repo.GetRoundByVoteId(request.Context(), voteId)
	if err != nil {
		writeInternalError(responseWriter, err)
		return
	}
	if round == nil {
		writeNotFound(responseWriter, "round not found")
		return
	}
	writeJSON(responseWriter, http.StatusOK, round)
}

// GET /api/Voting/GetResultForVote/?number={int}&partyId?&partyAcronim?
func (handler *VotingHandler) GetResultForVote(responseWriter http.ResponseWriter, request *http.Request) {
	numberStr := request.URL.Query().Get("number")
	if numberStr == "" {
		writeBadRequest(responseWriter, "missing 'number' parameter")
		return
	}

	voteId, err := strconv.Atoi(numberStr)
	if err != nil {
		writeBadRequest(responseWriter, "invalid 'number': must be an integer")
		return
	}

	voteFilter := repository.VoteFilter{
		RoundVoteId:  &voteId,
		PartyId:      request.URL.Query().Get("partyId"),
		PartyAcronim: request.URL.Query().Get("partyAcronim"),
	}

	votes, err := handler.repo.GetVotes(request.Context(), voteFilter)
	if err != nil {
		writeInternalError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, votes)
}

// GET /api/Voting/GetAllVotesForARoundById/?roundId={guid}
func (handler *VotingHandler) GetAllVotesForARoundById(responseWriter http.ResponseWriter, request *http.Request) {
	roundId := request.URL.Query().Get("roundId")
	if roundId == "" {
		writeBadRequest(responseWriter, "missing 'roundId' parameter")
		return
	}

	voteFilter := repository.VoteFilter{RoundID: roundId}

	votes, err := handler.repo.GetVotes(request.Context(), voteFilter)
	if err != nil {
		writeInternalError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, votes)
}

func collectKeywords(queryParams url.Values) []string {
	var keywords []string
	// ?keywords[]=foo form
	for _, keyword := range queryParams["keywords[]"] {
		if keyword != "" {
			keywords = append(keywords, keyword)
		}
	}
	// ?keywords=foo form (frontend sends this)
	for _, keyword := range queryParams["keywords"] {
		if keyword != "" {
			keywords = append(keywords, keyword)
		}
	}
	return keywords
}
