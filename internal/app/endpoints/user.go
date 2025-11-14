package endpoints

import (
	"net/http"

	"github.com/okunix/prservice/internal/pkg/models"
	"github.com/okunix/prservice/internal/pkg/pullrequest"
	"github.com/okunix/prservice/internal/pkg/user"
)

func SetIsActive(repo user.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req user.SetIsActiveRequest
		if err := ReadJson(r.Body, &req); err != nil {
			WriteJson(w, http.StatusBadRequest, models.ErrInvalidRequestBodyFormat)
			return
		}

		resp, err := repo.SetIsActive(r.Context(), req)
		if err != nil {
			if e, ok := err.(models.ErrorResponse); ok {
				statusCode := 400
				switch e.Err.Code {
				case models.NOT_FOUND:
					statusCode = 404
				}
				WriteJson(w, statusCode, e)
				return
			}
			WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		WriteJson(w, http.StatusOK, resp)
	}
}

func GetReview(repo pullrequest.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("user_id")
		resp, err := repo.GetReview(r.Context(), userId)
		if err != nil {
			WriteJson(w, http.StatusNotFound, err)
			return
		}
		WriteJson(w, http.StatusOK, resp)
	}
}
