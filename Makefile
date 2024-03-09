# build goodcommit
build:
	go fmt ./...
	go build -o bin/goodcommit cmd/main.go
