package backend

import "errors"

var (
	ErrUnauthorized           = errors.New("unauthorized")
	ErrNotFound               = errors.New("not found")
	ErrInternalServer         = errors.New("internal server error")
	ErrBadRequest             = errors.New("bad request")
	ErrConflict               = errors.New("conflict")
	ErrForbidden              = errors.New("forbidden")
	ErrServiceUnavailable     = errors.New("service unavailable")
	ErrTooManyRequests        = errors.New("too many requests")
	ErrInvalidInput           = errors.New("invalid input")
	ErrDatabaseError          = errors.New("database error")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrPermissionDenied       = errors.New("permission denied")
	ErrResourceNotFound       = errors.New("resource not found")
	ErrOperationNotAllowed    = errors.New("operation not allowed")
	ErrRateLimitExceeded      = errors.New("rate limit exceeded")
	ErrInvalidToken           = errors.New("invalid token")
	ErrSessionExpired         = errors.New("session expired")
	ErrInsufficientPrivileges = errors.New("insufficient privileges")
	ErrInvalidRequestFormat   = errors.New("invalid request format")
	ErrUnsupportedMediaType   = errors.New("unsupported media type")
	ErrMethodNotAllowed       = errors.New("method not allowed")
	ErrServiceTimeout         = errors.New("service timeout")
	ErrGatewayTimeout         = errors.New("gateway timeout")
	ErrConflictDetected       = errors.New("conflict detected")
	ErrResourceAlreadyExists  = errors.New("resource already exists")
	ErrInvalidState           = errors.New("invalid state")
)
