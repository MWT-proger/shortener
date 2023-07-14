package filestorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageGet(t *testing.T) {
	testCases := []struct {
		name        string
		tempStorage map[string]string
		key         string
		want        string
	}{
		{name: "Тест 1 - Проверяем на успех", tempStorage: map[string]string{"testKey": "testValue"}, key: "testKey", want: "testValue"},
		{name: "Тест 2 - Проверяем на пустую строку", tempStorage: map[string]string{"testKey": "testValue"}, key: "testKeyNot", want: ""},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := &FileStorage{
				tempStorage: tt.tempStorage,
			}
			got, _ := s.Get(tt.key)

			assert.Equal(t, tt.want, got, "Результат не совпадает с ожиданием")

		})
	}
}

func TestStorageSet(t *testing.T) {
	testCases := []struct {
		name        string
		tempStorage map[string]string
		value       string
		length      int
	}{
		{name: "Тест 1 - Проверяем на успех", tempStorage: map[string]string{"testKey0": "testValue0"}, value: "testValue", length: 2},
		{name: "Тест 2 - Проверяем на успех", tempStorage: map[string]string{"testKey0": "testValue0"}, value: "testValue2", length: 2},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := &FileStorage{
				tempStorage: tt.tempStorage,
			}
			got, _ := s.Set(tt.value)

			assert.Equal(t, tt.value, s.tempStorage[got], "Результат не совпадает с ожиданием")
			assert.Len(t, s.tempStorage, tt.length, "Длина словаря не совпадает")

		})
	}
}
