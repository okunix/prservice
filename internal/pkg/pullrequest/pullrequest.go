package pullrequest

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/okunix/prservice/internal/pkg/models"
)

type Reviewer struct {
	UserId        string `json:"user_id"`
	UserName      string `json:"user_name"`
	PullRequestId string `json:"pull_request_id"`
}

type Author struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PullRequest struct {
	Id                string     `json:"id"`
	Name              string     `json:"name"`
	Author            Author     `json:"author"`
	Status            Status     `json:"status"`
	Reviewers         []Reviewer `json:"reviewers"`
	NeedMoreReviewers bool       `json:"needMoreReviewers"`
}

type Status string

const (
	STATUS_OPEN   Status = "OPEN"
	STATUS_MERGED Status = "MERGED"
)

var (
	statusRegex = regexp.MustCompile(`^(?:OPEN|MERGED)$`)
	nameRegex   = regexp.MustCompile(`^[a-zA-Z0-9_]{1,100}$`)
)

var (
	ErrInvalidStatus    = errors.New("invalid status provided")
	ErrInvalidName      = errors.New("invalid name provided")
	ErrTooManyReviewers = errors.New("too many reviewers assigned")
	ErrAuthorReviewer   = errors.New("author of a pull request can't be it's reviewer")
)

func New(name string, author Author, initialReviewers []Reviewer) (*PullRequest, error) {
	if len(initialReviewers) > 2 {
		return nil, ErrTooManyReviewers
	}

	for _, v := range initialReviewers {
		if author.Id == v.UserId {
			return nil, ErrAuthorReviewer
		}
	}

	pr := PullRequest{
		Id:                uuid.NewString(),
		Name:              name,
		Author:            author,
		Reviewers:         initialReviewers,
		Status:            STATUS_OPEN,
		NeedMoreReviewers: false,
	}
	return &pr, pr.Validate()
}

func (pr *PullRequest) Validate() error {
	problems := models.ValidationError{}

	if err := validateName(pr.Name); err != nil {
		problems["name"] = err
	}

	if err := validateStatus(pr.Status); err != nil {
		problems["status"] = err
	}

	if len(problems) > 0 {
		return problems
	}
	return nil
}

func validateStatus(s Status) error {
	if !statusRegex.MatchString(string(s)) {
		return ErrInvalidStatus
	}
	return nil
}

func validateName(name string) error {
	if !nameRegex.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
