package save_test

import (
	"testing"
	mos "url_shortener/internal/http-server/handlers/url/mockers"
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mos.NewMockURLSaver

			if tt.respError == "" || tt.mockError != nil {
				urlSaverMock
			}
		})
	}
}
