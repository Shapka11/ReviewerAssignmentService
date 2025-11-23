package handler

const (
	ErrCodeBadRequest    = "BAD_REQUEST"
	ErrCodeInternalError = "INTERNAL_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeTeamExists    = "TEAM_EXISTS"
	ErrCodePRExists      = "PR_EXISTS"
	ErrCodePRMerged      = "PR_MERGED"
	ErrCodeNotAssigned   = "NOT_ASSIGNED"
	ErrCodeNoCandidate   = "NO_CANDIDATE"
)

const (
	ErrMsgInvalidJSON         = "invalid json body"
	ErrMsgAuthorNotFound      = "author not found"
	ErrMsgPRNotFound          = "pull request not found"
	ErrMsgPRMerged            = "cannot reassign on merged PR"
	ErrMsgPRExists            = "PR id already exists"
	ErrMsgReviewerNotAssigned = "reviewer is not assigned to this PR"
	ErrMsgNoCandidate         = "no active replacement candidate in team"
	ErrMsgTeamExists          = "team already exists"
	ErrMsgMissingTeamName     = "missing team_name"
	ErrMsgTeamNotFound        = "team not found"
	ErrMsgMissingUserID       = "missing user_id"
	ErrMsgUserNotFound        = "user not found"
)
