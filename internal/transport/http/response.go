package http

// errorResponse - структура для ответа с ошибкой.
type errorResponse struct {
	Message string `json:"message"`
}

// statusResponse - структура для ответа с сообщением о статусе.
type statusResponse struct {
	Status string `json:"status"`
}
