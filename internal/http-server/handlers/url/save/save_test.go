package save_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	_ "url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/logger/slogdiscard"

	"url-shortener/internal/http-server/handlers/url/save/mocks"

	"go.uber.org/mock/gomock"
)

func TestSaveHandler(t *testing.T) {

	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			alias:     "some_alias",
			url:       "",
			respError: "field URL is a required field",
		},
	}

	for _, tt := range tests {

		//чтобы каждый их этих тестов использовал свой нужный case
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			//тесты запускаются параллельно
			t.Parallel()

			ctr := gomock.NewController(t)

			defer ctr.Finish()
			urlSaverMock := mocks.NewMockURLSaver(ctr)

			// if tt.respError == "" || tt.mockError != nil {
			// 	urlSaverMock.
			// }

			handler := save.New(context.Background(), slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "#{tc.url}", "alias": "#{tc.alias}"}`)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			// gomock.N

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
		})
	}
}
