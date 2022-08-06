VERSION=develop
ifeq ($(OS),Windows_NT) 
    ext := '.exe'
else
    ext := ''
endif
build:
	@echo 'Building server'
	@cd src; \
	go build -ldflags="-X 'github.com/gabereiser/swr.version=$(VERSION)'" -o ../bin/server${ext};
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
