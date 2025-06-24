.PHONY: gen
gen:
	mockgen -source=internal/http-server/handlers/url/save/save.go \
	-destination=internal/http-server/handlers/url/mocks/mock_test.go


mockgen -source=C:\Users\office\Desktop\golang\url-shortener\internal\http-server\handlers\url\save\save.go -destination=C:\Users\office\Desktop\golang\url-shortener\internal\http-server\handlers\url\save\mocks\mock_test.go -package=mocks

работ


mockgen -source=./internal/http-server/handlers/url/save/save.go -destination=./internal/http-server/handlers/url/save/mocks/mock_save.go

mockgen -source=.\internal\http-server\handlers\url\save\save.go -destination=.\internal\http-server\handlers\url\save\mocks\mock_save.go

mockgen -source=internal/http-server/handlers/url/save/save.go -destination=internal/http-server/handlers/url/save/mocks/mock_save.go -package=mocks

