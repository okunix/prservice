package user

import "context"

type Repo interface {
	SetIsActive(ctx context.Context, req SetIsActiveRequest) (UserResponse, error)
	GetUsersByTeamName(ctx context.Context, teamName string) ([]User, error)
}
