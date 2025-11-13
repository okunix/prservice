package endpoints

import (
	"net/http"

	"github.com/okunix/prservice/internal/pkg/team"
)

func AddTeam(repo team.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req team.AddTeamRequest
		if err := ReadJson(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("can't parse request body"))
			return
		}

		resp, err := repo.AddTeam(r.Context(), req)
		if err != nil {
			WriteJson(w, http.StatusBadRequest, err)
			return
		}

		WriteJson(w, http.StatusCreated, resp)
	}
}
