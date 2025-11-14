package user

import "context"

type SetIsActiveRequest struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	User User `json:"user"`
}

type Repo interface {
	SetIsActive(ctx context.Context, req SetIsActiveRequest) (UserResponse, error)
	GetUsersByTeamName(ctx context.Context, teamName string) ([]User, error)
}
