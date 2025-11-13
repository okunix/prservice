package models

import "encoding/json"

type ErrorResponse struct {
	Err struct {
		Code    ErrorCode `json:"code"`
		Message string    `json:"message"`
	} `json:"error"`
}

type ErrorCode string

var (
	TEAM_EXISTS  ErrorCode = "TEAM_EXISTS"
	PR_EXISTS    ErrorCode = "PR_EXISTS"
	PR_MERGED    ErrorCode = "PR_MERGED"
	NOT_ASSIGNED ErrorCode = "NOT_ASSIGNED"
	NO_CANDIDATE ErrorCode = "NO_CANDIDATE"
	NOT_FOUND    ErrorCode = "NOT_FOUND"
)

func NewErrorResponse(code ErrorCode, msg string) ErrorResponse {
	return ErrorResponse{
		Err: struct {
			Code    ErrorCode "json:\"code\""
			Message string    "json:\"message\""
		}{
			Code:    code,
			Message: msg,
		},
	}
}

func (e ErrorResponse) Error() string {
	errBytes, _ := json.Marshal(e)
	return string(errBytes)
}

var (
	ErrTeamExists  = NewErrorResponse(TEAM_EXISTS, "team_name already exists")
	ErrNotFound    = NewErrorResponse(NOT_FOUND, "resource not found")
	ErrPRExists    = NewErrorResponse(PR_EXISTS, "PR is already exists")
	ErrPRMerged    = NewErrorResponse(PR_MERGED, "cannot reassign on merged PR")
	ErrNoCandidate = NewErrorResponse(NO_CANDIDATE, "no candidate")
	ErrNotAssigned = NewErrorResponse(NOT_ASSIGNED, "cannot assign user")
)
