package constants

// Context keys
type contextKey string

func (c contextKey) String() string {
	return "pkg/middleware context key " + string(c)
}

// Context key constants
const (
	ContextKeyUserID  = contextKey("uid")
	ContextKeyRoles   = contextKey("roles")
	ContextKeyTokenID = contextKey("tid")
)
