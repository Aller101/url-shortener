package save_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/http-server/handlers"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/lib/logger/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestValidURL(t *testing.T) {

	req := &save.Request{URL: "567944", Alias: "123"}

	fmt.Println(len(req.URL))
	fmt.Println(strings.TrimSpace(req.URL))
	fmt.Println(len(req.URL))
	fmt.Println(req.URL)
	vard := strings.TrimSpace(req.URL)
	fmt.Println(vard)
	fmt.Println(len(strings.TrimSpace(req.URL)))
	err := save.ValidReq(req)
	fmt.Println(len(req.URL))
	// if err != nil {
	// 	t.Fatalf("valid url must pass: %s", err.Error())
	// }
	require.NoError(t, err)

}

func TestValidErrURL(t *testing.T) {
	cases := []struct {
		name   string
		url    string
		expErr error
	}{
		{
			name:   "len = 4",
			url:    "1234",
			expErr: handlers.ErrLengtURL,
		},
		{
			name:   "len = 11",
			url:    "12345678910",
			expErr: handlers.ErrLengtURL,
		},
		{
			name:   "len = 0",
			url:    "",
			expErr: handlers.ErrVoidURL,
		},
		{
			name:   "len = 2x blank space",
			url:    "  ",
			expErr: handlers.ErrLengtURL,
		},
		{
			name:   "len = 8x blank space",
			url:    "        ",
			expErr: handlers.ErrLengtURL,
		},
		{
			name:   "len = 3 + 4x blank space",
			url:    "123    ",
			expErr: handlers.ErrLengtURL,
		},
		{
			name:   "len = 3x blank space + 5",
			url:    "   56789",
			expErr: handlers.ErrLengtURL,
		},
	}

	for _, tCase := range cases {
		req := save.Request{URL: tCase.url}
		t.Run(tCase.name, func(t *testing.T) {

			err := save.ValidReq(&req)
			require.EqualError(t, tCase.expErr, err.Error())
			require.Error(t, err)

		})
	}
}

func TestSaveURL(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockSaver := mocks.NewMockURLSaver(ctl)

	ctx := context.Background()
	// log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log := slogdiscard.NewDiscardLogger()

	cases := []struct {
		name         string
		requestBody  string
		mockSetup    func()
		expectedCode int
	}{
		{
			name:        "Success",
			requestBody: `{"url": "https://google.com", "alias": "test"}`,
			mockSetup: func() {
				// Ожидаем вызов SaveURL с конкретными параметрами
				mockSaver.EXPECT().
					SaveURL(gomock.Any(), "https://google.com", "test").
					Return(int64(1), nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Empty Alias - Auto Generate",
			requestBody: `{"url": "https://google.com"}`,
			mockSetup: func() {
				// Здесь нужно мокировать random.NewRandomString, если это возможно
				// Или ожидать любой alias
				mockSaver.EXPECT().
					SaveURL(gomock.Any(), "https://google.com", gomock.Any()).
					Return(int64(1), nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid Request",
			requestBody:  `{"url": ""}`,
			mockSetup:    func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "URL Already Exists",
			requestBody: `{"url": "https://exists.com", "alias": "exists"}`,
			mockSetup: func() {
				mockSaver.EXPECT().
					SaveURL(gomock.Any(), "https://exists.com", "exists").
					Return(int64(0), storage.ErrURLExists)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	tCase := cases[0]
	// var resp int64 = 1

	// input := fmt.Sprintf(`{"url": "#{tCase.url}", "alias": "#{tCase.alias}"}`)
	req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewBufferString(tCase.requestBody))

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	hh := save.New(ctx, log, mockSaver)

	//выполнили код
	hh(rr, req)

	//читаем
	res := rr.Result()

	_, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	// handler.ServeHTTP(nil, req)
	// require.NoError(t, err)

	// require.Equal(t, tCase.expectedCode, string(data))
	assert.Equal(t, tCase.expectedCode, rr.Code)

}
