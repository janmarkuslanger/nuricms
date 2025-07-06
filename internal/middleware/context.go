package middleware

type ctxKey string

const (
	UserIDKey    ctxKey = "userID"
	UserEmailKey ctxKey = "userEmail"
	UserRoleKey  ctxKey = "userRole"
)
