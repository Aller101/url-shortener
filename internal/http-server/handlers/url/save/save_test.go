package save_test

import (
	"fmt"
	"strings"
	"testing"
	"url-shortener/internal/http-server/handlers"
	"url-shortener/internal/http-server/handlers/url/save"

	"github.com/stretchr/testify/require"
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
