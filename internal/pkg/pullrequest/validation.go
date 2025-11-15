package pullrequest

import (
	"errors"
	"regexp"
	"slices"

	"github.com/okunix/prservice/internal/pkg/models"
)

var (
	statusRegex = regexp.MustCompile(`^(?:OPEN|MERGED)$`)
	nameRegex   = regexp.MustCompile(`^.{1,200}$`)
	idRegex     = regexp.MustCompile(`^[^ ]+$`)
)

var (
	ErrInvalidStatus    = errors.New("invalid status provided")
	ErrInvalidId        = errors.New("invalid id provided")
	ErrInvalidName      = errors.New("invalid name provided")
	ErrTooManyReviewers = errors.New("too many reviewers assigned")
	ErrAuthorReviewer   = errors.New("author of a pull request can't be it's reviewer")
)

func (pr *PullRequest) Validate() error {
	problems := models.ValidationError{}

	if err := ValidateId(pr.Id); err != nil {
		problems["id"] = err.Error()
	}

	if err := ValidateName(pr.Name); err != nil {
		problems["name"] = err.Error()
	}

	if err := ValidateStatus(pr.Status); err != nil {
		problems["status"] = err.Error()
	}

	if err := ValidateReviewers(pr.Reviewers, pr.AuthorId); err != nil {
		problems["reviewers"] = err.Error()
	}

	if len(problems) > 0 {
		return problems
	}

	return nil
}

func ValidateStatus(s Status) error {
	if !statusRegex.MatchString(string(s)) {
		return ErrInvalidStatus
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

func ValidateReviewers(reviewers []string, authorId string) error {
	if len(reviewers) > 2 {
		return ErrTooManyReviewers
	}
	if slices.Contains(reviewers, authorId) {
		return ErrAuthorReviewer
	}
	return nil
}
