install: tidy test
	go install ./cmd/skit/

tidy:
	go mod tidy

all: clean test build

test:
	go test -cover ./...

clean:
	rm -rf bin

build:
	mkdir -p bin/
	go build -o bin/skit cmd/skit/*.go

.PHONY: all	clean	build