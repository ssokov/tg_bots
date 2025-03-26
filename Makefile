.PHONY: build run

BINARY=mybot

build:
	go build -o $(BINARY) .


run: build
	./$(BINARY)
