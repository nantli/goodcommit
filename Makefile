# build goodcommit
build:
	go fmt ./...
	go build -o bin/goodcommit cmd/goodcommit/main.go
