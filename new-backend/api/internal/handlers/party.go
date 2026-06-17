package handlers

import (
	"net/http"
	"strconv"

	"github.com/borkanie/parliament-api/internal/repository"
)

type PartyHandler struct {
	repo *repository.PartyRepo
}

func NewPartyHandler(repo *repository.PartyRepo) *PartyHandler {
	return &PartyHandler{repo: repo}
}

// GET /api/Party/all?active={bool}&number={int=100}
func (handler *PartyHandler) All(responseWriter http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()

	var active *bool
	if activeStr := queryParams.Get("active"); activeStr != "" {
		parsedBool, err := strconv.ParseBool(activeStr)
		if err != nil {
			writeBadRequest(responseWriter, "invalid 'active' parameter")
			return
		}
		active = &parsedBool
	}

	limit := 100
	if numberStr := queryParams.Get("number"); numberStr != "" {
		count, err := strconv.Atoi(numberStr)
		if err != nil || count <= 0 {
			writeBadRequest(responseWriter, "invalid 'number' parameter")
			return
		}
		limit = count
	}

	parties, err := handler.repo.All(request.Context(), active, limit)
	if err != nil {
		writeInternalError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, parties)
}

// GET /api/Party/GetById/?id={guid}
func (handler *PartyHandler) GetByID(responseWriter http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		writeBadRequest(responseWriter, "missing 'id' parameter")
		return
	}

	party, err := handler.repo.GetByID(request.Context(), id)
	if err != nil {
		writeInternalError(responseWriter, err)
		return
	}
	if party == nil {
		writeNotFound(responseWriter, "party not found")
		return
	}
	writeJSON(responseWriter, http.StatusOK, party)
}

// GET /api/Party/query?name={str?}&acronym={str?}
func (handler *PartyHandler) Query(responseWriter http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()
	name := queryParams.Get("name")
	acronym := queryParams.Get("acronym")

	if name == "" && acronym == "" {
		writeBadRequest(responseWriter, "at least one of 'name' or 'acronym' is required")
		return
	}

	party, err := handler.repo.Query(request.Context(), name, acronym)
	if err != nil {
		writeInternalError(responseWriter, err)
		return
	}
	if party == nil {
		writeNotFound(responseWriter, "party not found")
		return
	}
	writeJSON(responseWriter, http.StatusOK, party)
}
