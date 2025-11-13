package team

import "context"

type Repo interface {
	AddTeam(ctx context.Context, t Team) (*Team, error)
	GetTeamByName(ctx context.Context, name string) (*Team, error)
}
