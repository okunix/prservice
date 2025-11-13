package pullrequest

import (
	"context"
	"time"
)

type PRMergeRequest struct {
	PullRequestId string `json:"pull_request_id"`
}

type PRCreateRequest struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
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

var (
	ErrResourceNotFound error
)

type Repo interface {
	// GET:/users/getReview
	GetReview(ctx context.Context, userId string) (GetReviewResponse, error)
	Merge(ctx context.Context, req PRMergeRequest) (PRMergeResponse, error)
	Reassign(ctx context.Context, req PRReassignRequest) (PRReassignResponse, error)
	Create(ctx context.Context, req PRCreateRequest) (PRCreateResponse, error)
	GetById(ctx context.Context, id string) (*PullRequest, error)
}
