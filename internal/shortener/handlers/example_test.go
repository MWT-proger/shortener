package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"

	"github.com/MWT-proger/shortener/internal/shortener/auth"
)

type MockTest struct{}

func (m *MockTest) Errorf(format string, args ...any) {}
func (m *MockTest) Fatalf(format string, args ...any) {}
func ExampleAPIHandler_GenerateShortkeyHandler() {

	var (
		ctrl        = gomock.NewController(&MockTest{})
		mockService = NewMockShortenerServicer(ctrl)

		h, _ = NewAPIHandler(mockService)

		bodyRequest = strings.NewReader("http://e.com")
		response    = httptest.NewRecorder()
		mockReq, _  = http.NewRequest(http.MethodPost, "http://localhost", bodyRequest)
		userID      = uuid.New()
		ctx         = auth.WithUserID(mockReq.Context(), userID)
	)
	mockReq = mockReq.WithContext(ctx)

	mockService.EXPECT().
		GenerateShortURL(mockReq.Context(), userID, "http://e.com", mockReq.Host).
		Return("http://localhost/testKey", http.StatusCreated)

	h.GenerateShortkeyHandler(response, mockReq)
	fmt.Println(response.Code)
	// Output:
	// 201

}
