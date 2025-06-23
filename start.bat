set CONFIG_PATH=./config/local.yaml
set GOOSE_DRIVER=postgres
set GOOSE_DBSTRING=postgresql://postgres:1233@localhost:5432/postgres?sslmode=disable

set CONFIG_PATH
set GOOSE_DRIVER
set GOOSE_DBSTRING

goose -dir migrations up

go run .\cmd\url-shortener\main.go