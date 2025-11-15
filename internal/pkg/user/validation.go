package user

import (
	"errors"
	"regexp"

	"github.com/okunix/prservice/internal/pkg/models"
)

var (
	idRegex   = regexp.MustCompile(`^[^ ]+$`)
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
)

var (
	ErrInvalidName = errors.New("invalid name provided")
	ErrInvalidId   = errors.New("invalid id provided")
)

func (u *User) Validate() error {
	problems := models.ValidationError{}

	if err := ValidateId(u.Id); err != nil {
		problems["id"] = err.Error()
	}
	if err := ValidateName(u.Name); err != nil {
		problems["name"] = err.Error()
	}

	if len(problems) > 0 {
		return problems
	}

	return nil
}

func ValidateName(name string) error {
	if !nameRegex.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}

func ValidateId(id string) error {
	if !idRegex.MatchString(id) {
		return ErrInvalidId
	}
	return nil
}
