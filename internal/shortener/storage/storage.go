package storage

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

const DB = "../../db.json"

func InitJSONFileStorage() {
	// Проверяет есть ли файл по указанному пути и если нет, создаёт его
	if _, err := os.ReadFile(DB); err != nil {
		str := "{}"
		if err = os.WriteFile(DB, []byte(str), 0644); err != nil {
			log.Fatal(err)
		}
	}
}

func SetInStorage(fullURL string) string {
	// Добавляет в хранилище полную ссылку и присваевает ей ключ
	shortURL := utils.StringWithCharset(5)

	dbJSON := make(map[string]string, 0)
	content, err := os.ReadFile(DB)

	if err == nil {
		if err = json.Unmarshal(content, &dbJSON); err == nil {

			for {
				_, ok := dbJSON[shortURL]
				if !ok {
					dbJSON[shortURL] = fullURL
					break
				}
				shortURL = utils.StringWithCharset(5)

			}

			b, err := json.Marshal(dbJSON)
			if err == nil {
				os.WriteFile(DB, b, 0644)
				return shortURL
			}
		}
	}
	return ""

}

func GetFromStorage(shortURL string) string {
	// Возвращает fullUrl по shortUrl

	dbJSON := make(map[string]string, 0)
	content, err := os.ReadFile(DB)

	if err == nil {
		if err = json.Unmarshal(content, &dbJSON); err == nil {

			fullURL, ok := dbJSON[shortURL]
			if ok {
				return fullURL
			}

		}
	}
	return ""

}
