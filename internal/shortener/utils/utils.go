package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/MWT-proger/shortener/configs"
)

const Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// StringWithCharset(length int) Выдаёт рандомную строку из указаннного колличества символов
func StringWithCharset(length int) string {
	// Временно использую для грубого теста
	// strong := []string{"sdsd", "pdunm"}
	// return strong[seededRand.Intn(len(strong))]

	b := make([]byte, length)
	for i := range b {
		b[i] = Charset[seededRand.Intn(len(Charset))]
	}
	return string(b)
}

// GetBaseShortURL(host string) формирует строку вида пример: http://localhost:8080
func GetBaseShortURL(host string) string {
	conf := configs.GetConfig()

	if conf.BaseURLShortener != "" {
		return fmt.Sprintf("%v/", conf.BaseURLShortener)
	} else {
		return fmt.Sprintf("http://%s/", host)
	}
}
