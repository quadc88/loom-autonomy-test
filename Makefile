.PHONY: all run test build clean docker-build docker-run

all: build

run:
	go run .

test:
	go test ./...

build:
	go build -o task-api .

docker-build:
	docker build -t task-api .

docker-run:
	docker run -p 8080:8080 task-api

clean:
	rm -f task-api
	rm -rf __pycache__
