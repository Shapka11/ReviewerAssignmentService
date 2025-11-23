package domains

type ReviewerStat struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	ReviewCount int    `json:"review_count"`
}

type GlobalStats struct {
	TotalUsers int `json:"total_users"`
	TotalPRs   int `json:"total_prs"`
}
