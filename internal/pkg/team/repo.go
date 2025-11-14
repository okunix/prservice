package team

import "context"

type AddTeamRequest struct {
	Team Team `json:"team"`
}

func (req *AddTeamRequest) Validate() error {
	return req.Team.Validate()
}

type AddTeamResponse struct {
	Team Team `json:"team"`
}

type GetTeamResponse Team

type Repo interface {
	AddTeam(ctx context.Context, t AddTeamRequest) (AddTeamResponse, error)
	GetTeamByName(ctx context.Context, name string) (GetTeamResponse, error)
}
