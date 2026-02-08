package handler

import (
	"errors"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/service"
)

// Возвращает http код и msg
func HTTPStatusError(err error) (int, string) {
	var svcErr *service.ServiceError
	if errors.As(err, &svcErr) {

		switch svcErr.Code {

		case service.ErrNotFound:
			return http.StatusNotFound, svcErr.Msg

		case service.ErrInvalidType:
			return http.StatusBadRequest, svcErr.Msg

		case service.ErrInvalidValue:
			return http.StatusBadRequest, svcErr.Msg

		default:
			return http.StatusInternalServerError, "internal error"

		}
	}

	return http.StatusInternalServerError, "internal error"
}
