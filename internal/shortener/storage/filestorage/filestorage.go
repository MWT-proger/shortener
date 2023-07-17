package filestorage

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

type FileStorage struct {
	storage.Storage
	tempStorage map[string]string
}

// InitJSONFile() Проверяет есть ли файл по указанному пути и если нет, создаёт его
func (s *FileStorage) Init(ctx context.Context) error {
	conf := configs.GetConfig()

	s.tempStorage = make(map[string]string, 0)
	if conf.JSONFileDB != "" {

		content, err := os.ReadFile(conf.JSONFileDB)

		if err != nil {
			str := "{}"
			if err = os.WriteFile(conf.JSONFileDB, []byte(str), 0644); err != nil {
				return err
			}

		} else {

			if err = json.Unmarshal(content, &s.tempStorage); err != nil {
				return err

			}
		}

		go s.BackupToJSONFile(ctx)

	}
	return nil

}

// BackupToJSONFile() Делает резервное копирование переменной в файл
func (s *FileStorage) BackupToJSONFile(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		default:
			time.Sleep(time.Minute * 10)
			log.Println("Старт Резервного копирования")

			conf := configs.GetConfig()

			b, err := json.Marshal(s.tempStorage)

			if err != nil {
				return err
			}

			os.WriteFile(conf.JSONFileDB, b, 0644)
			log.Println("Финиш Резервного копирования")

		}
	}

}

// Добавляет в хранилище полную ссылку и присваевает ей ключ
func (s *FileStorage) Set(fullURL string) (string, error) {

	shortURL := utils.StringWithCharset(5)

	for {
		_, ok := s.tempStorage[shortURL]
		if !ok {
			s.tempStorage[shortURL] = fullURL
			break
		}
		shortURL = utils.StringWithCharset(5)
	}

	return shortURL, nil

}

// Достаёт из хранилища и возвращает полную ссылку по ключу
func (s *FileStorage) Get(shortURL string) (string, error) {

	fullURL, ok := s.tempStorage[shortURL]
	if !ok {
		return "", nil
	}

	return fullURL, nil
}