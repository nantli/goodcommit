# build goodcommit
build:
	go fmt ./**/*.go
	go build -o bin/goodcommit cmd/main.go
