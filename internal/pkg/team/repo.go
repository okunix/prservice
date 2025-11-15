package team

import "context"

type Repo interface {
	AddTeam(ctx context.Context, t AddTeamRequest) (AddTeamResponse, error)
	GetTeamByName(ctx context.Context, name string) (GetTeamResponse, error)
	Deactivate(ctx context.Context, t DeactivateTeamRequest) (DeactivateTeamResponse, error)
}
