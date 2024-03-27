all: build

build:
	go build -o diode -ldflags="-X main.SemVer=0.1.0" diode.go

test:
	go test -v ./...

run:
	go run diode.go

