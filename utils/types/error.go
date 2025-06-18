package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v3"
)

type SuccessResponse struct {
	Data       interface{} `json:"data,omitempty"`
	Metadata   interface{} `json:"metadata,omitempty"`
	StatusCode int         `json:"-"`
}

type ErrorResponse struct {
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	Code       string      `json:"code"`
	StatusCode int         `json:"-"`
}

func NewSuccessResponse() *SuccessResponse {
	return &SuccessResponse{}
}

func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{}
}

func (s *SuccessResponse) WithStatusCode(code int) *SuccessResponse {
	s.StatusCode = code
	return s
}

func (e *ErrorResponse) WithStatusCode(code int) *ErrorResponse {
	e.StatusCode = code
	return e
}
func (e *ErrorResponse) WithCode(code string) *ErrorResponse {
	e.Code = code
	return e
}

func (e *ErrorResponse) WithErrorString(format string, a ...any) *ErrorResponse {
	e.Message = fmt.Sprintf(format, a...)
	return e
}

func (e *ErrorResponse) WithError(code error) *ErrorResponse {
	var valErrs validator.ValidationErrors
	if errors.As(code, &valErrs) {
		details := map[string]string{}
		for _, fieldError := range valErrs {
			fieldName := fieldError.Field()
			details["field"] = fieldName
			value, _ := json.Marshal(fieldError.Value())
			details["message"] = fmt.Sprintf("failed on tag %s, %s is not expected on %s", fieldError.Tag(), string(value), fieldError.ActualTag())
		}
		e.Details = details
		e.Code = "VALIDATION"
		e.Message = "Validation Error"
	} else if e.Code == "" {
		e.Code = "UNKNOWN"
	}

	e.Message = code.Error()
	return e
}
func (s *ErrorResponse) WithValidationError(fieldName string, message string) *ErrorResponse {
	s.Code = "VALIDATION"
	s.Message = "Validation Error"

	s.Details = map[string]string{
		"field":   fieldName,
		"message": message,
	}

	return s
}

func (s ErrorResponse) Error() string {
	return s.Message
}

func (s *SuccessResponse) WithData(data any) *SuccessResponse {
	s.Data = data
	return s
}
func (s *SuccessResponse) WithMetadata(Metadata any) *SuccessResponse {
	s.Metadata = Metadata
	return s
}
func (s SuccessResponse) Error() string {
	return "Success Response"
}

var (
	OKResponse       = NewSuccessResponse().WithStatusCode(http.StatusOK)
	AcceptedResponse = NewSuccessResponse().WithStatusCode(http.StatusAccepted)
	NoContentError   = NewErrorResponse().WithStatusCode(http.StatusNoContent)

	BadRequestError          = NewErrorResponse().WithCode("PARAMETERS").WithStatusCode(http.StatusBadRequest)
	UnauthorizedError        = NewErrorResponse().WithCode("UNAUTHORIZED").WithStatusCode(http.StatusUnauthorized)
	ForbiddenError           = NewErrorResponse().WithCode("FORBIDDEN").WithStatusCode(http.StatusForbidden)
	NotFoundError            = NewErrorResponse().WithStatusCode(http.StatusNotFound)
	UnprocessableEntityError = NewErrorResponse().WithCode("FETCH_DATA").WithStatusCode(http.StatusUnprocessableEntity)
	InternalServerError      = NewErrorResponse().WithCode("INTERNAL").WithStatusCode(http.StatusInternalServerError)
	NotImplementedError      = NewErrorResponse().WithCode("INTERNAL").WithStatusCode(http.StatusNotImplemented)
	BadGatewayError          = NewErrorResponse().WithCode("OVERLOAD").WithStatusCode(http.StatusBadGateway)
)

func CustomErrorHandler(c fiber.Ctx, err error) error {
	// Check if it's your custom Response error
	var resErr *SuccessResponse
	if errors.As(err, &resErr) {
		return c.Status(resErr.StatusCode).JSON(resErr.Data)
	}

	// Check if it's your custom Error error
	var errRes *ErrorResponse
	if errors.As(err, &errRes) {
		return c.Status(errRes.StatusCode).JSON(errRes)
	}

	// Fallback default error
	return fiber.DefaultErrorHandler(c, err)
}
