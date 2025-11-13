package repos

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/okunix/prservice/internal/pkg/models"
	"github.com/okunix/prservice/internal/pkg/pullrequest"
)

type PullRequestRepoImpl struct {
	db *sqlx.DB
}

func NewPullRequestRepo(db *sqlx.DB) pullrequest.Repo {
	return &PullRequestRepoImpl{
		db: db,
	}
}

func (p *PullRequestRepoImpl) Create(
	ctx context.Context,
	req pullrequest.PRCreateRequest,
) (pullrequest.PRCreateResponse, error) {
	var resp pullrequest.PRCreateResponse

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return resp, err
	}
	defer tx.Rollback()
	insertPullRequest := `
	INSERT INTO pull_requests (id, name, author_id, status) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, name, author_id, status;`
	err = tx.QueryRowContext(
		ctx,
		insertPullRequest,
		req.PullRequestId,
		req.PullRequestName,
		req.AuthorId,
		pullrequest.STATUS_OPEN,
	).Scan(
		&resp.PR.PullRequestId,
		&resp.PR.PullRequestName,
		&resp.PR.AuthorId,
		&resp.PR.Status,
	)
	if err != nil {
		slog.Error(err.Error())
		return resp, models.ErrPRExists
	}
	// also assigning 2 random users from team except author
	assignRandomUsersQuery := `
		INSERT INTO reviewers (user_id, pull_request_id) 
		SELECT id AS user_id, $1 FROM users 
		WHERE team_name = (SELECT u.team_name FROM pull_requests pr INNER JOIN users u ON u.id = pr.author_id WHERE pr.id = $1) 
		AND NOT id = (SELECT author_id FROM pull_requests WHERE id = $1) AND is_active = true ORDER BY RANDOM() LIMIT 2 RETURNING user_id;
	`
	rows, err := tx.QueryContext(ctx, assignRandomUsersQuery, resp.PR.PullRequestId)
	if err != nil {
		return resp, err
	}
	var assignedUsers []string
	for rows.Next() {
		var userId string
		if err := rows.Scan(&userId); err != nil {
			slog.Error(err.Error())
			return resp, err
		}
		assignedUsers = append(assignedUsers, userId)
	}
	resp.PR.AssignedReviewers = assignedUsers
	if len(assignedUsers) == 0 {
		updateNeedReviewersQuery := `UPDATE pull_requests SET need_more_reviewers = true WHERE id = $1;`
		_, err := tx.ExecContext(ctx, updateNeedReviewersQuery, resp.PR.PullRequestId)
		if err != nil {
			return resp, err
		}
	}
	return resp, tx.Commit()
}

func (p *PullRequestRepoImpl) GetReview(
	ctx context.Context,
	userId string,
) (pullrequest.GetReviewResponse, error) {
	selectUserPullRequests := `
		SELECT pr.id, pr.name, pr.author_id, pr.status 
		FROM pull_requests pr 
		INNER JOIN reviewers r ON pr.id = r.pull_request_id 
		WHERE r.user_id = $1;
	`
	var resp pullrequest.GetReviewResponse
	rows, err := p.db.QueryContext(ctx, selectUserPullRequests, userId)
	if err != nil {
		return resp, err
	}
	var prs []pullrequest.PullRequestShort
	for rows.Next() {
		var pr pullrequest.PullRequestShort
		if err := rows.Scan(&pr.Id, &pr.Name, &pr.AuthorId, &pr.Status); err != nil {
			return resp, err
		}
		prs = append(prs, pr)
	}
	resp.PullRequests = prs
	return resp, nil
}

func (p *PullRequestRepoImpl) Merge(
	ctx context.Context,
	req pullrequest.PRMergeRequest,
) (pullrequest.PRMergeResponse, error) {
	var resp pullrequest.PRMergeResponse

	pr, err := p.GetById(ctx, req.PullRequestId)
	if err != nil {
		return resp, err
	}
	if pr.Status == pullrequest.STATUS_MERGED {
		return resp, models.ErrPRMerged
	}

	mergedAt := time.Now()
	resp.PullRequestId = pr.Id
	resp.PullRequestName = pr.Name
	resp.AuthorId = pr.AuthorId
	resp.Status = pullrequest.STATUS_MERGED
	resp.MergedAt = mergedAt
	resp.AssignedReviewers = pr.Reviewers

	q := `UPDATE pull_requests SET status = $1, mergedAt = $2;`
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return resp, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, q, resp.Status, resp.MergedAt)
	if err != nil {
		return resp, err
	}

	return resp, tx.Commit()
}

func (p *PullRequestRepoImpl) Reassign(
	ctx context.Context,
	req pullrequest.PRReassignRequest,
) (pullrequest.PRReassignResponse, error) {
	var resp pullrequest.PRReassignResponse

	pr, err := p.GetById(ctx, req.PullRequestId)
	if err != nil {
		return resp, err
	}
	if pr.Status == pullrequest.STATUS_MERGED {
		return resp, models.ErrPRMerged
	}

	resp.PR = pullrequest.PR{
		PullRequestId:   pr.Id,
		PullRequestName: pr.Name,
		AuthorId:        pr.AuthorId,
		Status:          pr.Status,
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return resp, err
	}
	defer tx.Rollback()
	reassignQuery := `
		UPDATE reviewers SET user_id = (SELECT id AS user_id FROM users 
		WHERE team_name = (SELECT u.team_name FROM pull_requests pr INNER JOIN users u ON u.id = pr.author_id WHERE pr.id = $1) 
		AND NOT id = (SELECT author_id FROM pull_requests WHERE id = $1) 
		AND is_active = true AND NOT id = $2 ORDER BY RANDOM() LIMIT 1) WHERE pull_request_id = $1 AND user_id = $2 RETURNING user_id;
	`
	var replacedBy string
	err = tx.QueryRowContext(ctx, reassignQuery, req.PullRequestId, req.OldReviewerId).
		Scan(&replacedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resp, models.ErrNoCandidate
		}
		return resp, models.ErrNotAssigned
	}

	resp.ReplacedBy = replacedBy
	// adding new reviewer and removing old one
	resp.PR.AssignedReviewers = append(resp.PR.AssignedReviewers, resp.ReplacedBy)
	for i := 0; i < len(resp.PR.AssignedReviewers); i++ {
		if resp.PR.AssignedReviewers[i] == req.OldReviewerId {
			resp.PR.AssignedReviewers = append(
				resp.PR.AssignedReviewers[:i],
				resp.PR.AssignedReviewers[i+1:]...)
			break
		}
	}
	return resp, nil
}

func (p *PullRequestRepoImpl) GetById(
	ctx context.Context,
	id string,
) (*pullrequest.PullRequest, error) {
	q := `SELECT id, name, author_id, status, need_more_reviewers, created_at, merged_at FROM pull_requests WHERE id = $1;`
	var pr pullrequest.PullRequest
	err := p.db.QueryRowContext(ctx, q, id).
		Scan(
			&pr.Id,
			&pr.Name,
			&pr.AuthorId,
			&pr.Status,
			&pr.NeedMoreReviewers,
			&pr.CreatedAt,
			&pr.MergedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	getReviewersQuery := `SELECT user_id FROM reviewers WHERE pull_request_id = $1;`
	rows, err := p.db.QueryContext(ctx, getReviewersQuery, id)
	var reviewers []string
	for rows.Next() {
		var reviewer string
		if err := rows.Scan(&reviewer); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, reviewer)
	}
	pr.Reviewers = reviewers
	return &pr, nil
}
