package helper

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Field string
	Msg   string
}

func Error_handler(err error) []ApiError {
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ApiError, len(ve))
			for i, fe := range ve {
				out[i] = ApiError{fe.Field(), msgForTag(fe.Tag())}
			}
			return out
		}
		return nil
	}
	return nil
}

func Single_Error_handler(err error) ApiError {
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				return ApiError{fe.Field(), msgForTag(fe.Tag())}
			}

		}
		return ApiError{}
	}
	return ApiError{}
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "oneof":
		return "Invalid params"
	}
	return ""
}
