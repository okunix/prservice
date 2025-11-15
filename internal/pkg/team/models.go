package team

import (
	"github.com/okunix/prservice/internal/pkg/models"
	"github.com/okunix/prservice/internal/pkg/user"
)

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

type DeactivateTeamResponse struct {
	Team Team `json:"team"`
}

type DeactivateTeamRequest struct {
	TeamName string `json:"team_name"`
}

type TeamMember struct {
	Id       string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (m *TeamMember) Validate() error {
	problems := models.ValidationError{}
	if err := user.ValidateId(m.Id); err != nil {
		problems["user_id"] = err.Error()
	}
	if err := user.ValidateName(m.Name); err != nil {
		problems["username"] = err.Error()
	}
	if len(problems) > 0 {
		return problems
	}
	return nil
}

type Team struct {
	Name    string       `json:"team_name"`
	Members []TeamMember `json:"members"`
}

func (t *Team) Validate() error {
	problems := models.ValidationError{}

	if err := ValidateName(t.Name); err != nil {
		problems["name"] = err.Error()
	}

	if len(problems) > 0 {
		return problems
	}

	for _, v := range t.Members {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}
