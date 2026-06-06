BINARY=rigctld-mqtt
PREFIX=/usr/local/bin

all: build

build:
	mkdir -p bin
	go build -o bin/${BINARY} cmd/daemon/main.go

install: build
	install -d ${PREFIX}
	install -m 755 bin/${BINARY} ${PREFIX}/${BINARY}

run:
	go run cmd/daemon/main.go
