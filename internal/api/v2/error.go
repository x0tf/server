package v2

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// apiError represents an error returned from the REST API
type apiError struct {
	StatusCode int                    `json:"-"`
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
}

// Error formats the API error into a string
func (err *apiError) Error() string {
	return fmt.Sprintf("%d: %s (custom_data: %v)", err.Code, err.Message, err.Data)
}

// newError creates a new API error and makes sure the data map is initialized even when it is a nil value
func newError(statusCode, code int, message string, data map[string]interface{}) *apiError {
	if data == nil {
		data = map[string]interface{}{}
	}
	return &apiError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Data:       data,
	}
}

// ErrorHandler represents the fiber error handler to return API errors
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var apiError *apiError
	if !errors.As(err, &apiError) {
		apiError = newError(fiber.StatusInternalServerError, errorCodeUnexpectedError, err.Error(), nil)
		fmt.Println(apiError.StatusCode)
	}
	return ctx.Status(apiError.StatusCode).JSON(apiError)
}
