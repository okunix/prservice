package user

import "context"

type Repo interface {
	SetIsActive(ctx context.Context, userId string, isActive bool) (*User, error)
	GetUsersByTeamName(ctx context.Context, teamName string) ([]User, error)
}
