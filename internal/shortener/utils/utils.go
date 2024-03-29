package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/MWT-proger/shortener/configs"
)

// const charset - набор символов для генерации строки.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// seededRand Rand that uses random values from src to generate other random values.
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// StringWithCharset(length int) Выдаёт рандомную строку из указаннного колличества символов.
func StringWithCharset(length int) string {

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GetBaseShortURL(host string) формирует строку.
// пример: http://localhost:8080/
func GetBaseShortURL(conf configs.Config, host string) string {

	if conf.BaseURLShortener != "" {
		return fmt.Sprintf("%v/", conf.BaseURLShortener)
	} else {
		return fmt.Sprintf("http://%s/", host)
	}
}
