AUTHOR="kudohamu"
GOVERSION=$(shell go version)
PROJECT_NAME="spirali"
TARGETS=$(addprefix github.com/$(AUTHOR)/$(PROJECT_NAME)/cmd/, $(PROJECT_NAME))
TARGET_FILES=$(shell go list ./... 2> /dev/null | grep -v "/vendor/")

all: $(TARGETS)

$(TARGETS):
	@go install $@

.PHONY: lint vet test

lint:
	@lint=`golint $(shell ls -d | grep -v "vendor") 2>&1`; \
		lint=`echo "$$lint" | grep -E -v -e vendor/.+\.go`; \
		echo "$$lint"; if [ "$$lint" != "" ]; then exit 1; fi

vet:
	@vet=`go tool vet -all -structtags -tests $(shell ls -d */ | grep -v "vendor") 2>&1`; \
		vet=`echo "$$vet" | grep -E -v -e vendor/.+\.go`; \
		echo "$$vet"; if [ "$$vet" != "" ]; then exit 1; fi

test:
	@go test -v $(TARGET_FILES)
