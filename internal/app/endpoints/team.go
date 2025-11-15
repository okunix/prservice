package endpoints

import (
	"net/http"

	"github.com/okunix/prservice/internal/pkg/models"
	"github.com/okunix/prservice/internal/pkg/team"
)

func AddTeam(repo team.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req team.AddTeamRequest
		if err := ReadAndValidate(r.Body, &req); err != nil {
			WriteJson(w, http.StatusBadRequest, err)
			return
		}

		resp, err := repo.AddTeam(r.Context(), req)
		if err != nil {
			if validationError, ok := err.(models.ValidationError); ok {
				WriteJson(w, http.StatusBadRequest,
					models.ErrValidationFailed(validationError))
				return
			}
			WriteJson(w, http.StatusBadRequest, err)
			return
		}

		WriteJson(w, http.StatusCreated, resp)
	}
}

func GetTeamByName(repo team.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")
		resp, err := repo.GetTeamByName(r.Context(), teamName)
		if err != nil {
			WriteJson(w, http.StatusNotFound, err)
			return
		}
		WriteJson(w, http.StatusOK, resp)
	}
}

func Deactivate(repo team.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req team.DeactivateTeamRequest
		if err := ReadJson(r.Body, &req); err != nil {
			WriteJson(w, http.StatusBadRequest,
				models.ErrInvalidRequestBodyFormat)
			return
		}
		resp, err := repo.Deactivate(r.Context(), req)
		if err != nil {
			WriteJson(w, http.StatusNotFound, models.ErrNotFound)
			return
		}
		WriteJson(w, http.StatusOK, resp)
	}
}
