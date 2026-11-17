package api_runner

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	maxRequestAttempts = 8
	retryBaseDelay     = 500 * time.Millisecond
	retryMaxDelay      = 3 * time.Second
	httpClientTimeout  = 30 * time.Second
)

// повторяет запрос, если сервер временно ответил 503
type retryTransport struct {
	base http.RoundTripper
}

// отправляет запрос и при ответе 503 повторяет его с нарастающей паузой
func (t retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		var err error
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		_ = req.Body.Close()
	}

	for attempt := 1; ; attempt++ {
		attemptReq := req.Clone(req.Context())
		if body != nil {
			attemptReq.Body = io.NopCloser(bytes.NewReader(body))
		}

		resp, err := t.base.RoundTrip(attemptReq)
		if err != nil || resp.StatusCode != http.StatusServiceUnavailable || attempt >= maxRequestAttempts {
			return resp, err
		}

		_ = resp.Body.Close()
		delay := min(retryBaseDelay*time.Duration(attempt), retryMaxDelay)
		log.Printf("Сервер ответил 503, повторяем запрос (попытка %d из %d, пауза %s)",
			attempt, maxRequestAttempts, delay)
		time.Sleep(delay)
	}
}

// создаёт HTTP-клиент для тестов с автоматическим повтором запросов при 503
func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: retryTransport{base: http.DefaultTransport},
		Timeout:   httpClientTimeout,
	}
}
