build:
	go build -o bin/go_shell main.go
run:
	go run main.go

compile:		
	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go

install:
	go install 

all: build install
