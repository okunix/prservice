package user

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/okunix/prservice/internal/pkg/models"
)

type User struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	IsActive bool    `json:"is_active"`
	TeamName *string `json:"teamName,omitempty"`
}

func New(name string, isActive bool) (*User, error) {
	u := User{
		Id:       uuid.NewString(),
		Name:     name,
		IsActive: isActive,
		TeamName: nil,
	}
	return &u, u.Validate()
}

func (u *User) AssignTeam(teamName string) {
	u.TeamName = &teamName
}

func (u *User) Validate() error {
	problems := models.ValidationError{}

	if err := validateId(u.Id); err != nil {
		problems["id"] = err
	}
	if err := validateName(u.Name); err != nil {
		problems["name"] = err
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
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
)

func validateName(name string) error {
	if !nameRegex.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}

func validateId(id string) error {
	return uuid.Validate(id)
}
