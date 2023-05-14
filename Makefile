.PHONY: build run clean

build:
	go build -o ReadFile ./cmd

run:
	go run ./cmd

clean:
	rm -f ReadFile
