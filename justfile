check: build fmt

build:
	go build main.go
	rm main

fmt:
	go fmt ./...

test run=".*":
	go test ./... -run={{ run }}
