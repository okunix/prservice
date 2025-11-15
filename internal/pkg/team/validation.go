package team

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidName = errors.New("invalid name provided")
)

var (
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\- ]{1,40}$`)
)

func ValidateName(name string) error {
	if !nameRegex.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
