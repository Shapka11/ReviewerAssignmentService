package postgres

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ReviewerAssignmentService/internal/domains"
)

type teamRepositoryImpl struct {
	database *pgxpool.Pool
}

func NewTeamRepository(database *pgxpool.Pool) *teamRepositoryImpl {
	return &teamRepositoryImpl{database: database}
}

func (t *teamRepositoryImpl) Create(ctx context.Context, team *domains.Team) error {
	tx, err := t.database.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("warning: transaction rollback failed: %v", err)
		}
	}()

	queryTeam := `
        WITH ins AS (
            INSERT INTO teams (team_name) 
            VALUES ($1)
            ON CONFLICT (team_name) DO NOTHING
            RETURNING id
        )
        SELECT id FROM ins
        UNION ALL
        SELECT id FROM teams WHERE team_name = $1
        LIMIT 1;
    `
	var teamID int
	if err := tx.QueryRow(ctx, queryTeam, team.Name).Scan(&teamID); err != nil {
		return err
	}

	queryUser := `
		INSERT INTO users (user_id, username, team_id, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE 
		SET username = EXCLUDED.username, 
		    team_id = EXCLUDED.team_id,
		    is_active = EXCLUDED.is_active
	`

	for _, member := range team.Members {
		_, err = tx.Exec(ctx, queryUser, member.UserID, member.UserName, teamID, member.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (t *teamRepositoryImpl) Exists(ctx context.Context, teamName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`
	var exists bool

	err := t.database.QueryRow(ctx, query, teamName).Scan(&exists)

	return exists, err
}

func (t *teamRepositoryImpl) GetByName(ctx context.Context, name string) (*domains.Team, error) {
	var teamID int
	err := t.database.QueryRow(ctx, "SELECT id FROM teams WHERE team_name = $1", name).Scan(&teamID)
	if err != nil {
		return nil, nil
	}

	memberQuery := `
		SELECT user_id, username, is_active
		FROM users 
		WHERE team_id = $1
	`

	rows, err := t.database.Query(ctx, memberQuery, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]domains.TeamMember, 0)
	for rows.Next() {
		var member domains.TeamMember
		err = rows.Scan(&member.UserID, &member.UserName, &member.IsActive)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return &domains.Team{
		Name:    name,
		Members: members,
	}, nil
}
