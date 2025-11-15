package pullrequest

import (
	"time"

	"github.com/okunix/prservice/internal/pkg/models"
)

type PRMergeRequest struct {
	PullRequestId string `json:"pull_request_id"`
}

type PRCreateRequest struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

func (req PRCreateRequest) Validate() error {
	problems := models.ValidationError{}

	if err := ValidateName(req.PullRequestName); err != nil {
		problems["pull_request_name"] = err.Error()
	}

	if err := ValidateId(req.PullRequestId); err != nil {
		problems["pull_request_id"] = err.Error()
	}

	if len(problems) > 0 {
		return problems
	}
	return nil
}

type PRReassignRequest struct {
	PullRequestId string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}

type PRMergeResponse struct {
	PullRequestId     string    `json:"pull_request_id"`
	PullRequestName   string    `json:"pull_request_name"`
	AuthorId          string    `json:"author_id"`
	Status            Status    `json:"status"`
	AssignedReviewers []string  `json:"assigned_reviewers"`
	MergedAt          time.Time `json:"mergedAt"`
}

type PR struct {
	PullRequestId     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorId          string   `json:"author_id"`
	Status            Status   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

type PullRequestShort struct {
	Id       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
	Status   Status `json:"status"`
}

type GetReviewResponse struct {
	UserId       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

type PRReassignResponse struct {
	PR         PR     `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}

type PRCreateResponse struct {
	PR PR `json:"pr"`
}

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
