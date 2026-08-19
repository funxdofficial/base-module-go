PKG ?=
TYPE ?= rest
OUTPUT ?=

CLI := ./bin/base-module-go

help:
	@echo "Targets: tidy, test, build, install, fun-build, generate-example, generate-consumer-example, integration-test"
	@echo ""
	@echo "Generate via Make:"
	@echo "  make fun-build PKG=order-service TYPE=rest"
	@echo "  make fun-build PKG=sync-worker TYPE=consumer [OUTPUT=/path]"

tidy:
	@go mod tidy

test: tidy
	@go test ./...

build: tidy
	@mkdir -p bin
	@go build -o $(CLI) .

install:
	@go install .

integration-test: tidy
	@BASE_MODULE_GO_INTEGRATION=1 go test ./... -run TestGenerateServiceBuild -count=1

fun-build: build
	@test -n "$(PKG)" || (echo "Usage: make fun-build PKG=<service-name> [TYPE=rest|consumer|cons] [OUTPUT=<dir>]" && exit 1)
	@$(CLI) -pkg=$(PKG) -type=$(TYPE) $(if $(OUTPUT),-output=$(OUTPUT),)

generate-example:
	@$(MAKE) fun-build PKG=demo-service TYPE=rest OUTPUT=/tmp/demo-service
	@echo "Generated REST service at /tmp/demo-service"

generate-consumer-example:
	@$(MAKE) fun-build PKG=demo-worker TYPE=consumer OUTPUT=/tmp/demo-worker
	@echo "Generated consumer service at /tmp/demo-worker"
