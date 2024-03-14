all: build

build:
	go build -o diode -ldflags="-X main.SemVer=0.0.8" diode.go

test:
	go test -v

run:
	go run diode.go

