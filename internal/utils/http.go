package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// DecodeBody читает тело ответа и разбирает JSON в структуру указанного типа
func DecodeBody[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result T
	require.NoError(t, json.Unmarshal(body, &result), "Не удалось разобрать тело ответа: %s", string(body))
	return result
}
