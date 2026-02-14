package constants

// User roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User status
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusBanned   = "banned"
)

// Context keys
const (
	ContextKeyUserID    = "user_id"
	ContextKeyUserEmail = "user_email"
	ContextKeyUserRole  = "user_role"
	ContextKeyRequestID = "request_id"
)

// Header keys
const (
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
	HeaderRequestID     = "X-Request-ID"
	HeaderUserAgent     = "User-Agent"
)

// Cache keys
const (
	CacheKeyUserPrefix    = "user:"
	CacheKeyTokenPrefix   = "token:"
	CacheKeySessionPrefix = "session:"
)

// Cache TTL
const (
	CacheTTLShort  = 300  // 5 minutes
	CacheTTLMedium = 1800 // 30 minutes
	CacheTTLLong   = 3600 // 1 hour
)

// Queue names
const (
	QueueUserEvents = "user.events"
	QueueEmailQueue = "email.queue"
)

// Exchange names
const (
	ExchangeUserEvents = "user.events"
)

// Routing keys
const (
	RoutingKeyUserCreated = "user.created"
	RoutingKeyUserUpdated = "user.updated"
	RoutingKeyUserDeleted = "user.deleted"
)
