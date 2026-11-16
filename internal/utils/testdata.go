package utils

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Ptr возвращает указатель на переданное значение (для необязательных полей запросов)
func Ptr[T any](value T) *T {
	return &value
}

// UniqueTitle возвращает название с уникальным суффиксом, чтобы поиск находил только наше объявление
func UniqueTitle(prefix string) string {
	return fmt.Sprintf("%s %d", prefix, time.Now().UnixMilli())
}

// CreateTestPNG создаёт маленький PNG-файл во временной папке теста и возвращает путь к нему
func CreateTestPNG(t *testing.T, fileName string) string {
	t.Helper()

	filePath := filepath.Join(t.TempDir(), fileName)
	file, err := os.Create(filePath)
	require.NoError(t, err)
	defer file.Close()

	require.NoError(t, png.Encode(file, image.NewRGBA(image.Rect(0, 0, 10, 10))))
	return filePath
}
