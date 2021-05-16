package entity

import (
	"encoding/json"
	"time"
)

// ErrorResponse represents common error object
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents common success object
type SuccessResponse struct {
	Status string `json:"status"`
}

// NewError creates new object of ErrorResponse
func NewError(s string) ErrorResponse {
	return ErrorResponse{Error: s}
}

// NewErrorJSON creates new object of ErrorResponse and return byte array
func NewErrorJSON(s string) ([]byte, error) {
	e := ErrorResponse{Error: s}
	return json.Marshal(e)
}

// RemoveEmptyValues remove empty entries from map
func RemoveEmptyValues(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if v == nil {
			delete(m, k)
		}

		if v != nil {
			// handle empty string
			if tmp, ok := v.(string); ok && len(tmp) == 0 {
				delete(m, k)
			}

			// handle time
			if tmp, ok := v.(time.Time); ok && tmp.IsZero() {
				delete(m, k)
			}
		}
	}
	return m
}
