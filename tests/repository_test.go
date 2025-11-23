package tests

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ReviewerAssignmentService/internal/domains"
	internalPostgres "ReviewerAssignmentService/internal/repository/postgres"
)

func setupDB(t *testing.T) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), "postgres://user:password@localhost:5435/testdb")
	if err != nil || pool.Ping(context.Background()) != nil {
		t.Skip("Database not available")
	}

	_, err = pool.Exec(context.Background(), "TRUNCATE TABLE pull_requests, users, teams CASCADE")
	require.NoError(t, err, "failed to truncate tables")
	return pool
}

func TestRepository_TeamAndUser(t *testing.T) {
	pool := setupDB(t)
	defer pool.Close()

	team := &domains.Team{
		Name: "RepoTestTeam",
		Members: []domains.TeamMember{
			{UserID: "rt_u1", UserName: "User1", IsActive: true},
		},
	}

	err := internalPostgres.NewTeamRepository(pool).Create(context.Background(), team)
	require.NoError(t, err)

	exists, _ := internalPostgres.NewTeamRepository(pool).Exists(context.Background(), "RepoTestTeam")
	assert.True(t, exists)

	user, _ := internalPostgres.NewUserRepository(pool).GetByID(context.Background(), "rt_u1")
	assert.Equal(t, "RepoTestTeam", user.TeamName)
}

func TestRepository_PullRequest(t *testing.T) {
	pool := setupDB(t)
	defer pool.Close()

	ctx := context.Background()
	_, err := pool.Exec(ctx, "INSERT INTO teams (team_name) VALUES ('PRTeam') ON CONFLICT DO NOTHING")
	require.NoError(t, err, "failed to seed team")

	var teamID int
	err = pool.QueryRow(ctx, "SELECT id FROM teams WHERE team_name='PRTeam'").Scan(&teamID)
	require.NoError(t, err, "failed to get team id")

	_, err = pool.Exec(ctx, "INSERT INTO users (user_id, username, team_id) VALUES ('pr_author', 'Auth', $1)", teamID)
	require.NoError(t, err, "failed to seed user")

	pr := &domains.PullRequest{
		ID:                "pr-repo-1",
		Name:              "Test Repo",
		AuthorID:          "pr_author",
		Status:            domains.PRStatusOpen,
		AssignedReviewers: []string{"u2", "u3"},
	}

	err = internalPostgres.NewPrRepository(pool).Create(ctx, pr)
	require.NoError(t, err)

	fetched, _ := internalPostgres.NewPrRepository(pool).GetByID(ctx, "pr-repo-1")
	assert.Equal(t, "pr-repo-1", fetched.ID)
}
