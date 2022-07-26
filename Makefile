
build:
	@echo 'Building server'
	go build -o ../bin/server;
	@echo 'Done'
clean:
	@echo 'Cleaning build'
	@cd src; \
	go clean; \
	rm -rf ../bin/*
	@echo 'Done'
dependencies:
	@echo 'Downloading dependencies'
	@cd src; \
	go mod tidy; \
	go mod download
	@echo 'Done'

all: clean dependencies build
