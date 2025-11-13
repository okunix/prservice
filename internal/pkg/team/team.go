package team

import (
	"errors"
	"regexp"

	"github.com/okunix/prservice/internal/pkg/models"
	"github.com/okunix/prservice/internal/pkg/user"
)

type TeamMember struct {
	Id       string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	Name    string       `json:"team_name"`
	Members []TeamMember `json:"members"`
}

func New(name string) (*Team, error) {
	t := Team{
		Name: name,
	}
	return &t, t.Validate()
}

func (t *Team) AddMember(member user.User) {
	t.Members = append(t.Members, TeamMember{Id: member.Id, Name: member.Name})
}

func (t *Team) Validate() error {
	problems := models.ValidationError{}

	if err := validateName(t.Name); err != nil {
		problems["name"] = err.Error()
	}

	if len(problems) > 0 {
		return problems
	}
	return nil
}

var (
	ErrInvalidName = errors.New("invalid name provided")
)

var (
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{1,40}$`)
)

func validateName(name string) error {
	if !nameRegex.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
