package utils

type message struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Message(success bool, message string) map[string]interface{} {
	return map[string]interface{}{"success": success, "message": message}
}
