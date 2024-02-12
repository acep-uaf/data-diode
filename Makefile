build:
	go build -o diode -ldflags="-X main.SemVer=0.0.6" main.go

test:
	go test -v

run:
	go run main.go

