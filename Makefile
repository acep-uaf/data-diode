BIN_NAME=diode
BIN_VERSION=0.2.1
BIN_DATE=$(shell date +%FT%T%z)

all: build

build:
	go build -o ${BIN_NAME} -ldflags="-X 'main.SemVer=${BIN_VERSION}' -X 'main.BuildInfo=${BIN_DATE}'"

test:
	go test -v ./...

coverage:
	go test -cover ./...

run: build
	./${BIN_NAME} --help

clean:
	go clean
	rm ${BIN_NAME}

