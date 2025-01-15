package helpers

func APIResponse(message string, code int, status string, data interface{}) map[string]interface{} {
	meta := map[string]interface{}{
		"message": message,
		"code":    code,
		"status":  status,
	}

	if data == nil {
		data = map[string]interface{}{}
	}
	jsonResponse := map[string]interface{}{
		"meta": meta,
		"data": data,
	}

	return jsonResponse
}
