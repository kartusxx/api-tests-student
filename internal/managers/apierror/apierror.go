package apierror

// ErrorResponse описывает тело ответа API с ошибкой
type ErrorResponse struct {
	Details string `json:"details"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
