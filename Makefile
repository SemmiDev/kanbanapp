run-app:
	go run main.go
build-app:
	go build -o kanbanapp main.go
run-build-app:
	./kanbanapp
test-app:
	go test -v ./...
tailwind:
	npx tailwindcss -i ./input.css -o ./static/css/style.css --watch
proto:
	rm -f internal/pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=internal/proto --go_out=internal/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=kanbanapp \
	internal/proto/*.proto

.PHONY: run-app test-app tailwind proto build-app run-build-app
