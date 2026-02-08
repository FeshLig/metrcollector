package service

type ErrorCode string

const (
	ErrNotFound     ErrorCode = "not_found"
	ErrInvalidType  ErrorCode = "invalid_type"
	ErrInvalidValue ErrorCode = "invalid_value"
)

type ServiceError struct {
	Code ErrorCode
	Msg  string
}

func (e *ServiceError) Error() string {
	return e.Msg
}
