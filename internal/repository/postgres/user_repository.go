package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ReviewerAssignmentService/internal/domains"
)

type userRepositoryImpl struct {
	database *pgxpool.Pool
}

func NewUserRepository(database *pgxpool.Pool) *userRepositoryImpl {
	return &userRepositoryImpl{database: database}
}

func (u *userRepositoryImpl) Create(ctx context.Context, user *domains.User) error {
	query := `
        INSERT INTO users (user_id, username, team_id, is_active)
        VALUES ($1, $2, (SELECT id FROM teams WHERE team_name = $3), $4)
		ON CONFLICT (user_id) DO UPDATE
		SET is_active = EXCLUDED.is_active
    `

	_, err := u.database.Exec(ctx, query,
		user.ID, user.Name, user.TeamName, user.IsActive)
	return err
}

func (u *userRepositoryImpl) Exists(ctx context.Context, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`
	var exists bool
	err := u.database.QueryRow(ctx, query, userID).Scan(&exists)
	return exists, err
}

func (u *userRepositoryImpl) GetByID(ctx context.Context, id string) (*domains.User, error) {
	var user domains.User
	query := `
		SELECT u.user_id, u.username, t.team_name, u.is_active
		FROM users u
		JOIN teams t ON u.team_id = t.id
		WHERE u.user_id = $1
	`

	err := u.database.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.TeamName,
		&user.IsActive,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (u *userRepositoryImpl) GetRandomActiveUsersByTeam(
	ctx context.Context, teamName string, excludeUserID string, limit int) ([]string, error) {

	query := `
        SELECT u.user_id 
        FROM users u
        JOIN teams t ON u.team_id = t.id
        WHERE t.team_name = $1 
          AND u.is_active = TRUE 
          AND u.user_id != $2
        ORDER BY RANDOM()
        LIMIT $3
    `

	rows, err := u.database.Query(ctx, query, teamName, excludeUserID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userIDs := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}

func (u *userRepositoryImpl) UpdateActivity(ctx context.Context, userID string, isActive bool) error {
	query := `UPDATE users SET is_active = $2 WHERE user_id = $1`
	tag, err := u.database.Exec(ctx, query, userID, isActive)
	if err == nil && tag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return err
}

func (u *userRepositoryImpl) DeactivateTeamMembers(ctx context.Context, teamName string) error {
	query := `
		UPDATE users 
		SET is_active = false 
		WHERE team_id = (SELECT id FROM teams WHERE team_name = $1)
	`
	_, err := u.database.Exec(ctx, query, teamName)
	return err
}

func (u *userRepositoryImpl) Count(ctx context.Context) (int, error) {
	var count int
	err := u.database.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}
