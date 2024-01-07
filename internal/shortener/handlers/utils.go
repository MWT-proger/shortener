package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
)

func (h *APIHandler) unmarshalBody(body io.ReadCloser, form interface{}) error {

	defer body.Close()

	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf.Bytes(), form); err != nil {
		return err
	}

	return nil
}

// setHTTPError(w http.ResponseWriter, err error) присваивает response статус ответа
// вынесен для исключения дублирования в коде
func (h *APIHandler) setHTTPError(w http.ResponseWriter, err error) {
	var serviceError *lErrors.ServicesError
	if errors.As(err, &serviceError) {
		http.Error(w, serviceError.Error(), serviceError.HTTPCode)
	} else {
		http.Error(w, "Ошибка сервера, попробуйте позже.", http.StatusInternalServerError)
	}
}

func (h *APIHandler) setOrGetHTTPCode(w http.ResponseWriter, err error) int {
	var serviceError *lErrors.ServicesError

	if errors.As(err, &serviceError) {

		if !serviceError.IsReturn {
			return serviceError.HTTPCode
		} else {
			http.Error(w, serviceError.Error(), serviceError.HTTPCode)
			return 0
		}

	} else {
		http.Error(w, "Ошибка сервера, попробуйте позже.", http.StatusInternalServerError)
		return 0
	}
}
