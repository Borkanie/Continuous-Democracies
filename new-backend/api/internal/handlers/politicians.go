package handlers

import (
	"net/http"
	"strconv"

	"github.com/borkanie/parliament-api/internal/repository"
)

type PoliticiansHandler struct {
	repo *repository.PoliticiansRepo
}

func NewPoliticiansHandler(repo *repository.PoliticiansRepo) *PoliticiansHandler {
	return &PoliticiansHandler{repo: repo}
}

// GET /api/Politicians/getAllPoliticians
func (handler *PoliticiansHandler) GetAll(responseWriter http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()

	filter := repository.PoliticianFilter{
		PartyAcronym: queryParams.Get("partyAcronym"),
		PartyName:    queryParams.Get("partyName"),
		Limit:        100,
	}

	if activeStr := queryParams.Get("isActive"); activeStr != "" {
		parsedBool, err := strconv.ParseBool(activeStr)
		if err != nil {
			writeBadRequest(responseWriter, "invalid 'isActive' parameter")
			return
		}
		filter.IsActive = &parsedBool
	}

	if locationStr := queryParams.Get("location"); locationStr != "" {
		locationNum, err := strconv.Atoi(locationStr)
		if err != nil {
			writeBadRequest(responseWriter, "invalid 'location' parameter")
			return
		}
		filter.Location = &locationNum
	}

	if genderStr := queryParams.Get("gender"); genderStr != "" {
		genderNum, err := strconv.Atoi(genderStr)
		if err != nil {
			writeBadRequest(responseWriter, "invalid 'gender' parameter")
			return
		}
		filter.Gender = &genderNum
	}

	if numberStr := queryParams.Get("number"); numberStr != "" {
		count, err := strconv.Atoi(numberStr)
		if err != nil || count <= 0 {
			writeBadRequest(responseWriter, "invalid 'number' parameter")
			return
		}
		filter.Limit = count
	}

	politicians, err := handler.repo.GetAll(request.Context(), filter)
	if err != nil {
		writeInternalError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, politicians)
}

// GET /api/Politicians/GetById/?id={guid}
func (handler *PoliticiansHandler) GetByID(responseWriter http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		writeBadRequest(responseWriter, "missing 'id' parameter")
		return
	}

	politician, err := handler.repo.GetByID(request.Context(), id)
	if err != nil {
		writeInternalError(responseWriter, err)
		return
	}
	if politician == nil {
		writeNotFound(responseWriter, "politician not found")
		return
	}
	writeJSON(responseWriter, http.StatusOK, politician)
}

// GET /api/Politicians/GetByName/?name={str}
func (handler *PoliticiansHandler) GetByName(responseWriter http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	if name == "" {
		writeBadRequest(responseWriter, "missing 'name' parameter")
		return
	}

	politician, err := handler.repo.GetByName(request.Context(), name)
	if err != nil {
		writeInternalError(responseWriter, err)
		return
	}
	if politician == nil {
		writeNotFound(responseWriter, "politician not found")
		return
	}
	writeJSON(responseWriter, http.StatusOK, politician)
}
