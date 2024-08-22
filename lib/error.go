package lib

import (
	"bytes"
	"fmt"

)

type APIError struct {
	// Human readable message which corresponds to the client error
	// eg: "Invalid user"
	Message string `json:"message"`

	// Underscored delimited string
	// eg: "invalid"
	Code string `json:"code"`

	// Http Status Code from uber
	Status int `json:"status"`

	// A hash of field names that have validations. This has a value of an array with
	// member strings that describe the specific validation error
	// eg: map{"first_name": ["Required"]}
	Fields map[string]string `json:"fields,omitempty"`

	MessageError string `json:"message_error,omitempty"`

	// Metadata Error
	Meta interface{} `json:"meta,omitempty"`
}


func (err APIError) Error() string {
	var apiErrBuff bytes.Buffer

	if err.MessageError != "" {
		return err.MessageError
	}

	apiErrBuff.WriteString(fmt.Sprintf("API: %s", err.Message))

	if err.Code != "" {
		apiErrBuff.WriteString(fmt.Sprintf("\nCode: %s", err.Code))
	}

	if err.Fields != nil {
		apiErrBuff.WriteString("\nFields:")
		for k, v := range err.Fields {
			apiErrBuff.WriteString(fmt.Sprintf("\n\t%s: %v", k, v))
		}
	}

	return apiErrBuff.String()
}


func NewAPIError(message, code, messageError string, status int) APIError {
	return APIError{
		Message:      message,
		Code:         code,
		Status:       status,
		MessageError: messageError,
	}
}
