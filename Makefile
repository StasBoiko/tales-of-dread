.PHONY: all build frontend test lint run clean

all: build

frontend:
	npm --prefix frontend run build

build: frontend
	go build -o bin/server ./cmd/server

test:
	go test -v -race ./...

lint:
	go vet ./...

run: build
	./bin/server

clean:
	rm -rf bin/ web/dist/
