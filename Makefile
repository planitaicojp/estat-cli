.PHONY: build test lint clean

build:
	go build -o estat .

test:
	go test ./... -v

lint:
	golangci-lint run ./...

vet:
	go vet ./...

clean:
	rm -f estat
	rm -rf dist/
