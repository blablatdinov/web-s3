fmt:
	go fmt ./...

run:
	go run src/cmd/server/main.go

build:
	go build -o bin/server src/cmd/server/main.go
