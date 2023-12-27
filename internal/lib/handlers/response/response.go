package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type BaseResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	statusOk    = "OK"
	statusError = "Error"
)

func Ok() BaseResponse {
	return BaseResponse{Status: statusOk}
}

func Error(msg string) BaseResponse {
	return BaseResponse{Status: statusError, Error: msg}
}

func ValidationError(errs validator.ValidationErrors) BaseResponse {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s is required", err.Field()))
		}
	}
	return BaseResponse{
		Status: statusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
