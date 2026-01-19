package dto

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"мessage"`
	Details interface{} `json:"details,omitempty"`
}

func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

func NewErrorResponceWithDetailse(code, message string, details interface{}) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}
