build:
	go build -o diode testbed.go utility.go

test:
	go test -v

run:
	go run testbed.go utility.go

