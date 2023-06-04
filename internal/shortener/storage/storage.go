package storage

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

const DB = "../../db.json"

func InitJsonFileStorage() {
	// Проверяет есть ли файл по указанному пути и если нет, создаёт его
	if _, err := os.ReadFile(DB); err != nil {
		str := "{}"
		if err = os.WriteFile(DB, []byte(str), 0644); err != nil {
			log.Fatal(err)
		}
	}
}

func SetInStorage(fullUrl string) string {
	// Добавляет в хранилище полную ссылку и присваевает ей ключ
	shortUrl := utils.StringWithCharset(5)

	dbJson := make(map[string]string, 0)
	content, err := os.ReadFile(DB)

	if err == nil {
		if err = json.Unmarshal(content, &dbJson); err == nil {

			for {
				_, ok := dbJson[shortUrl]
				if !ok {
					dbJson[shortUrl] = fullUrl
					break
				}
				shortUrl = utils.StringWithCharset(5)

			}

			b, err := json.Marshal(dbJson)
			if err == nil {
				os.WriteFile(DB, b, 0644)
				return shortUrl
			}
		}
	}
	return ""

}

func GetFromStorage(shortUrl string) string {
	// Возвращает fullUrl по shortUrl

	dbJson := make(map[string]string, 0)
	content, err := os.ReadFile(DB)

	if err == nil {
		if err = json.Unmarshal(content, &dbJson); err == nil {

			fullUrl, ok := dbJson[shortUrl]
			if ok {
				return fullUrl
			}

		}
	}
	return ""

}
