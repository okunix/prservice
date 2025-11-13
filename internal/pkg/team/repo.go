package team

import "context"

type AddTeamRequest struct {
	Team Team `json:"team"`
}

type Repo interface {
	AddTeam(ctx context.Context, t AddTeamRequest) (*Team, error)
	GetTeamByName(ctx context.Context, name string) (*Team, error)
}
