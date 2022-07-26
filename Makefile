
build:
	$(eval GIT_COMMIT=0.0.5-$(shell git rev-parse --short head))
	@echo 'Building version $(GIT_COMMIT)'
	@cd src; \
	sed 's/$$VERSION/$(GIT_COMMIT)/g' ./swr/version.go.inc > ./swr/version.go; \
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
