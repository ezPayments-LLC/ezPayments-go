package ezpayments

import "fmt"

// APIError represents an error returned by the ezPayments API.
// Use errors.As to check for API errors:
//
//	var apiErr *ezpayments.APIError
//	if errors.As(err, &apiErr) {
//	    fmt.Println(apiErr.Code, apiErr.Message)
//	}
type APIError struct {
	// Type is the error type (e.g. "invalid_request_error", "authentication_error").
	Type string `json:"type"`

	// Message is a human-readable error description.
	Message string `json:"message"`

	// Code is a machine-readable error code (e.g. "resource_missing").
	Code string `json:"code"`

	// Param identifies the parameter related to the error, if applicable.
	Param string `json:"param"`

	// StatusCode is the HTTP status code of the response.
	StatusCode int `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("ezpayments: %s (type=%s, code=%s, status=%d)", e.Message, e.Type, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("ezpayments: %s (type=%s, status=%d)", e.Message, e.Type, e.StatusCode)
}

// apiErrorEnvelope is the JSON structure for API error responses.
type apiErrorEnvelope struct {
	Error APIError `json:"error"`
}
