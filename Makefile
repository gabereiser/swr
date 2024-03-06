VERSION=develop
ifeq ($(OS),Windows_NT) 
    ext := '.exe'
else
    ext := ''
endif

build:
	@echo 'Building server'
	@go build -ldflags="-X 'github.com/gabereiser/swr.version=$(VERSION)'" -o ./bin/server${ext};
	@echo 'Done'
clean:
	@echo 'Cleaning build'
	@go clean; \
	rm -rf ./bin; \
	rm -rf ./release;
	@echo 'Done'
dependencies:
	@echo 'Downloading dependencies'
	@go mod tidy; \
	go mod download
	@echo 'Done'
rel:
	@echo 'Making release'
	@cp -R data bin/data; \
	cp -R docs bin/docs; \
	mkdir -p release; \
	rm -f release/swr-$(VERSION).tar.gz; \
	tar -czf release/swr-$(VERSION).tar.gz -C bin .
	@echo 'Done releasing version $(VERSION)'
all: clean dependencies build rel
