package storage

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

const DB = "../../db.json"

type OperationStorage interface {
	SetInStorage(fullURL string) (string, error)
	GetFromStorage(shortURL string) (string, error)
}
type Storage struct {
}

func InitJSONFileStorage() {
	// Проверяет есть ли файл по указанному пути и если нет, создаёт его
	if _, err := os.ReadFile(DB); err != nil {
		str := "{}"
		if err = os.WriteFile(DB, []byte(str), 0644); err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Storage) SetInStorage(fullURL string) (string, error) {
	// Добавляет в хранилище полную ссылку и присваевает ей ключ
	shortURL := utils.StringWithCharset(5)

	dbJSON := make(map[string]string, 0)
	content, err := os.ReadFile(DB)

	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(content, &dbJSON); err != nil {
		return "", err
	}

	for {
		_, ok := dbJSON[shortURL]
		if !ok {
			dbJSON[shortURL] = fullURL
			break
		}
		shortURL = utils.StringWithCharset(5)
	}

	b, err := json.Marshal(dbJSON)
	if err != nil {
		return "", err
	}

	os.WriteFile(DB, b, 0644)
	return shortURL, nil

}

func (s *Storage) GetFromStorage(shortURL string) (string, error) {
	// Достаёт из хранилища и возвращает полную ссылку по ключу

	dbJSON := make(map[string]string, 0)
	content, err := os.ReadFile(DB)

	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(content, &dbJSON); err != nil {
		return "", err
	}

	fullURL, ok := dbJSON[shortURL]
	if !ok {
		return "", nil
	}

	return fullURL, nil
}
