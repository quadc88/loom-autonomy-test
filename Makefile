.PHONY: all run test build clean

all: build

run:
	python app.py

test:
	npm test

build:
	docker build -t task-api .

clean:
	rm -f .coverage
	rm -rf __pycache__
	rm -rf node_modules
	find . -name '*.pyc' -delete
