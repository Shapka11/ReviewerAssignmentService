package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ReviewerAssignmentService/internal/domains"
)

type prRepositoryImpl struct {
	database *pgxpool.Pool
}

func NewPrRepository(database *pgxpool.Pool) *prRepositoryImpl {
	return &prRepositoryImpl{database: database}
}

func (p *prRepositoryImpl) Create(ctx context.Context, pr *domains.PullRequest) error {
	query := `
        INSERT INTO pull_requests (
            pull_request_id,
            pull_request_name,
            author_id,
            status,
            assigned_reviewers
        ) VALUES ($1, $2, $3, $4, $5)
    `

	reviewersJSON, err := json.Marshal(pr.AssignedReviewers)
	if err != nil {
		return err
	}

	_, err = p.database.Exec(ctx, query,
		pr.ID,
		pr.Name,
		pr.AuthorID,
		pr.Status,
		string(reviewersJSON),
	)
	return err
}

func (p *prRepositoryImpl) Exists(ctx context.Context, prID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`
	var exists bool
	err := p.database.QueryRow(ctx, query, prID).Scan(&exists)
	return exists, err
}

func (p *prRepositoryImpl) GetByID(ctx context.Context, id string) (*domains.PullRequest, error) {
	query := `
		SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	var pr domains.PullRequest
	var reviewersJSON []byte

	row := p.database.QueryRow(ctx, query, id)
	err := row.Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&reviewersJSON,
		&pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(reviewersJSON) > 0 {
		if err := json.Unmarshal(reviewersJSON, &pr.AssignedReviewers); err != nil {
			return nil, err
		}
	} else {
		pr.AssignedReviewers = []string{}
	}

	return &pr, nil
}

func (p *prRepositoryImpl) GetByReviewer(ctx context.Context, reviewerID string) ([]*domains.PullRequestShort, error) {
	query := `
        SELECT 
            pull_request_id,
            pull_request_name,
            author_id,
            status
        FROM pull_requests 
        WHERE assigned_reviewers @> to_jsonb($1::text)
        ORDER BY created_at DESC
    `

	rows, err := p.database.Query(ctx, query, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*domains.PullRequestShort
	for rows.Next() {
		var pr domains.PullRequestShort
		err = rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.AuthorID,
			&pr.Status,
		)
		if err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}

	return prs, nil
}

func (p *prRepositoryImpl) Update(ctx context.Context, pr *domains.PullRequest) error {
	query := `
        UPDATE pull_requests 
        SET pull_request_name = $2,
            status = $3::varchar,
            assigned_reviewers = $4,
            merged_at = CASE WHEN $3::varchar = 'MERGED' AND merged_at IS NULL THEN CURRENT_TIMESTAMP ELSE merged_at END
        WHERE pull_request_id = $1
    `

	reviewersJSON, err := json.Marshal(pr.AssignedReviewers)
	if err != nil {
		return err
	}

	_, err = p.database.Exec(ctx, query,
		pr.ID, pr.Name, pr.Status, string(reviewersJSON))

	return err
}

func (p *prRepositoryImpl) Count(ctx context.Context) (int, error) {
	var count int
	err := p.database.QueryRow(ctx, "SELECT COUNT(*) FROM pull_requests").Scan(&count)
	return count, err
}
