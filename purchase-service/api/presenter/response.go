package presenter

// ErrorResponse creates a standardized error response
func ErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"error":   message,
	}
}

// SuccessResponse creates a standardized success response
func SuccessResponse(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"data":    data,
	}
}
