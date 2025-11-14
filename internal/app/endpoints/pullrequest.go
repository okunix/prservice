package endpoints

import (
	"net/http"

	"github.com/okunix/prservice/internal/pkg/models"
	"github.com/okunix/prservice/internal/pkg/pullrequest"
)

func CreatePullRequest(repo pullrequest.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pullrequest.PRCreateRequest
		if err := ReadJson(r.Body, &req); err != nil {
			WriteJson(w, http.StatusBadRequest,
				models.ErrInvalidRequestBodyFormat)
			return
		}
		resp, err := repo.Create(r.Context(), req)
		if err != nil {
			WriteJson(w, http.StatusConflict, err)
			return
		}

		WriteJson(w, http.StatusCreated, resp)
	}
}

func MergePullRequest(repo pullrequest.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pullrequest.PRMergeRequest
		if err := ReadJson(r.Body, &req); err != nil {
			WriteJson(w, http.StatusBadRequest,
				models.ErrInvalidRequestBodyFormat)
			return
		}
		resp, err := repo.Merge(r.Context(), req)
		if err != nil {
			WriteJson(w, http.StatusBadRequest, err)
			return
		}
		WriteJson(w, http.StatusOK, resp)
	}
}

func ReassignPullRequest(repo pullrequest.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pullrequest.PRReassignRequest
		if err := ReadJson(r.Body, &req); err != nil {
			WriteJson(w, http.StatusBadRequest,
				models.ErrInvalidRequestBodyFormat)
			return
		}
		resp, err := repo.Reassign(r.Context(), req)
		if err != nil {
			WriteJson(w, http.StatusBadRequest, err)
			return
		}
		WriteJson(w, http.StatusOK, resp)
	}
}
