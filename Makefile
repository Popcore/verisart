EXECUTABLE=verisart
BUILD_DIR=build

.PHONY: build
build:
	@echo "==> building executable"
	@mkdir -p build/
	go build -o $(BUILD_DIR)/$(EXECUTABLE)

.PHONY: test
test:
	@echo "==> running unit tests"
	go test ./...