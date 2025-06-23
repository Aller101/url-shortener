go install github.com/pressly/goose/v3/cmd/goose@latest


go install github.com/golang/mock/mockgen@v1.6.0

mockgen -source=internal/http-server/handlers/url/save/save.go -destination=internal/http-server/handlers/url/mock/mock_test.go