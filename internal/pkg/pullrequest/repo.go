package pullrequest

import (
	"context"
)

type Repo interface {
	// GET:/users/getReview
	GetReview(ctx context.Context, userId string) (GetReviewResponse, error)
	Merge(ctx context.Context, req PRMergeRequest) (PRMergeResponse, error)
	Reassign(ctx context.Context, req PRReassignRequest) (PRReassignResponse, error)
	Create(ctx context.Context, req PRCreateRequest) (PRCreateResponse, error)
	GetById(ctx context.Context, id string) (*PullRequest, error)
}
