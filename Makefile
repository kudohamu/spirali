AUTHOR="kudohamu"
GOVERSION=$(shell go version)
PROJECT_NAME="spirali"
TARGETS=$(addprefix github.com/$(AUTHOR)/$(PROJECT_NAME)/cmd/, $(PROJECT_NAME))
TARGET_FILES=$(shell go list ./... 2> /dev/null | grep -v "/vendor/")

all: $(TARGETS)

$(TARGETS):
	@go install $@

.PHONY: test

test:
	@go test -v $(TARGET_FILES)
