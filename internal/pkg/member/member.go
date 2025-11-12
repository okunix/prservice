package member

import (
	"github.com/okunix/prservice/internal/pkg/team"
	"github.com/okunix/prservice/internal/pkg/user"
)

type TeamMember struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	TeamName string `json:"team_name"`
}

func New(user user.User, team team.Team) TeamMember {
	return TeamMember{
		UserId:   user.Id,
		UserName: user.Name,
		TeamName: team.Name,
	}
}
