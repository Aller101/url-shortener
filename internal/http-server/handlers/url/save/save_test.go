package save_test

import (
	"testing"
	_ "url-shortener/internal/http-server/handlers/url/save"

	// "url-shortener/internal/http-server/handlers/url/mocks"
	"url-shortener/internal/http-server/handlers/url/save/mocks"

	"go.uber.org/mock/gomock"
)

func TestSaveHandler(t *testing.T) {
	gomock.NewController(t)
	mocks.NewMockURLSaver()

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mos.NewMockURLSaver

			if tt.respError == "" || tt.mockError != nil {
				urlSaverMock
			}
		})
	}
}
