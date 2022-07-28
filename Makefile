VERSION=develop

build:
	@echo 'Building server'
	@cd src; \
	go build -ldflags="-X 'github.com/gabereiser/swr.version=$(VERSION)'" -o ../bin/server;
	@echo 'Done'
clean:
	@echo 'Cleaning build'
	@cd src; \
	go clean; \
	rm -rf ../bin/*
	@echo 'Done'
client:
	@echo 'Building client'
	@cd src/client; \
	go build -ldflags="-X 'main.version=$(VERSION)'" -o ../../bin/client;
	@echo 'Done'
dependencies:
	@echo 'Downloading dependencies'
	@cd src; \
	go mod tidy; \
	go mod download
	@echo 'Done'

all: clean dependencies build client
