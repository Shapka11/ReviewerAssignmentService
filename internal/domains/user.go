package domains

type User struct {
	ID       string `json:"user_id" db:"user_id"`
	Name     string `json:"username" db:"username"`
	TeamName string `json:"team_name" db:"team_name"`
	IsActive bool   `json:"is_active" db:"is_active"`
}
