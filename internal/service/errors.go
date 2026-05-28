package service

// ErrorCode represents service-level error type.
type ErrorCode string

const (
	// ErrNotFound indicates that requested metric was not found.
	ErrNotFound ErrorCode = "not_found"

	// ErrInvalidType indicates unsupported metric type.
	ErrInvalidType ErrorCode = "invalid_type"

	// ErrInvalidValue indicates invalid metric value.
	ErrInvalidValue ErrorCode = "invalid_value"
)

// ServiceError describes business logic error returned by service layer.
type ServiceError struct {
	Code ErrorCode
	Msg  string
}

// Error returns service error message.
func (e *ServiceError) Error() string {
	return e.Msg
}
