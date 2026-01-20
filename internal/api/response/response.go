package response

import (
	"encoding/json"
	"net/http"
)

// Response is the standard API response envelope.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo contains error details.
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// Meta contains pagination and other metadata.
type Meta struct {
	Total  int64 `json:"total,omitempty"`
	Limit  int   `json:"limit,omitempty"`
	Offset int   `json:"offset,omitempty"`
}

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Response{
		Success: status >= 200 && status < 300,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

// JSONWithMeta writes a JSON response with pagination metadata.
func JSONWithMeta(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Response{
		Success: status >= 200 && status < 300,
		Data:    data,
		Meta:    meta,
	}

	json.NewEncoder(w).Encode(resp)
}

// Error writes an error response.
func Error(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}

	json.NewEncoder(w).Encode(resp)
}

// ValidationError writes a validation error response.
func ValidationError(w http.ResponseWriter, field, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	resp := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "validation_error",
			Message: message,
			Field:   field,
		},
	}

	json.NewEncoder(w).Encode(resp)
}

// NotFound writes a not found response.
func NotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "resource not found"
	}
	Error(w, http.StatusNotFound, "not_found", message)
}

// Unauthorized writes an unauthorized response.
func Unauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "unauthorized"
	}
	Error(w, http.StatusUnauthorized, "unauthorized", message)
}

// Forbidden writes a forbidden response.
func Forbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "forbidden"
	}
	Error(w, http.StatusForbidden, "forbidden", message)
}

// InternalError writes an internal server error response.
func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal_error", "an internal error occurred")
}

// BadRequest writes a bad request response.
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "bad_request", message)
}
