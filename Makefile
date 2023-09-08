all: test

test:
	go test -v .

deps:
	go get -u github.com/go-ble/ble
	go get -u github.com/raff/goble
	go get -u github.com/nats-io/nats.go

README.md:
	go test .
	go run github.com/campoy/embedmd@latest -w README.md

.PHONY:all test deps README.md
