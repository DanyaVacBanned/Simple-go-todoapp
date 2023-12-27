BINARY_NAME=todo-app
.DEFAULT_GOAL := run

dep:
	go mod install

build:
	go build  ./cmd/main/main.go

run: build
	./build/${BINARY_NAME}
