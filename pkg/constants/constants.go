package constants

// User roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Pagination defaults
const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// Context keys
const (
	ContextKeyUserID    = "userID"
	ContextKeyUserEmail = "userEmail"
	ContextKeyUserRole  = "userRole"
)

// Time formats
const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
	TimeFormat     = "15:04:05"
)
