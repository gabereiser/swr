VERSION=develop
ifeq ($(OS),Windows_NT) 
    ext := '.exe'
else
    ext := ''
endif
GO_MAJOR_VERSION = $(shell go version|cut -c 14-|cut -d. -f1)
GO_MIDDLE_VERSION = $(shell go version|cut -d. -f2)
GO_MINIMUM_MAJOR_VERSION = $(shell grep ^go go.mod|cut -d" " -f2|cut -d. -f1)
GO_MINIMUM_MIDDLE_VERSION = $(shell grep ^go go.mod|cut -d" " -f2|cut -d. -f2)
define check_go_version
    @if [ $(($GO_MAJOR_VERSION * 1000 + $GO_MIDDLE_VERSION)) -lt $(($GO_MINIMUM_MAJOR_VERSION * 1000 + $GO_MINIMUM_MIDDLE_VERSION)) ]; then \
	echo 'You need an updated version of go: you are currently using $(GO_MAJOR_VERSION).$(GO_MIDDLE_VERSION) but you need $(GO_MINIMUM_MAJOR_VERSION).$(GO_MINIMUM_MIDDLE_VERSION) or later.'; \
	exit 1; \
    fi
endef

build:
	@echo 'Building server'
	$(call check_go_version)
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
	$(call check_go_version)
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
