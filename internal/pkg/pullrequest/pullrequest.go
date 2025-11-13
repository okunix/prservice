package pullrequest

import (
	"errors"
	"regexp"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/okunix/prservice/internal/pkg/models"
)

type PullRequest struct {
	Id                string     `json:"pull_request_id"`
	Name              string     `json:"pull_request_name"`
	AuthorId          string     `json:"author_id"`
	Status            Status     `json:"status"`
	Reviewers         []string   `json:"assigned_reviewers"`
	NeedMoreReviewers bool       `json:"needMoreReviewers"`
	CreatedAt         time.Time  `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
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

func New(name string, authorId string, reviewerIds []string) (*PullRequest, error) {
	if len(reviewerIds) > 2 {
		return nil, ErrTooManyReviewers
	}

	if slices.Contains(reviewerIds, authorId) {
		return nil, ErrAuthorReviewer
	}

	pr := PullRequest{
		Id:                uuid.NewString(),
		Name:              name,
		AuthorId:          authorId,
		Reviewers:         reviewerIds,
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
