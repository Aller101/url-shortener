package s_test

import (
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
)

func TestValid(t *testing.T) {

	req := &save.Request{URL: "123456", Alias: "123"}

	err := save.ValidReq(req)
	if err != nil {
		t.Fatalf("valid url must pass: %s", err.Error())
	}

}
