.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: gen
gen:
	mockgen -source=internal/http-server/handlers/url/save/save.go \
	-destination=internal/http-server/handlers/url/mocks/mock_save.go
