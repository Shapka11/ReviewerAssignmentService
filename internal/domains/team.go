package domains

type TeamMember struct {
	UserID   string `json:"user_id" db:"user_id"`
	UserName string `json:"username" db:"username"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type Team struct {
	Name    string       `json:"team_name" db:"team_name"`
	Members []TeamMember `json:"members" db:"members"`
}
