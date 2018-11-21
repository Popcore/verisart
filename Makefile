EXECUTABLE=verisart
BUILD_DIR=build
ARTEFACT_DIR=artefacts

.PHONY: build
build:
	@echo "==> building executable"
	@mkdir -p build/
	go build -o $(BUILD_DIR)/$(EXECUTABLE)

.PHONY: test
test:
	@echo "==> running unit tests"
	@mkdir -p $(ARTEFACT_DIR)
	@echo 'mode: atomic' > $(ARTEFACT_DIR)/coverage.out
	go test ./... -coverprofile=$(ARTEFACT_DIR)/coverage.tmp && tail -n +2 $(ARTEFACT_DIR)/coverage.tmp >> $(ARTEFACT_DIR)/coverage.out || exit;
	go tool cover -html=$(ARTEFACT_DIR)/coverage.out -o $(ARTEFACT_DIR)/coverage.html