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
	@go test ./... -coverprofile=$(ARTEFACT_DIR)/coverage.out
	@go tool cover -html=$(ARTEFACT_DIR)/coverage.out -o $(ARTEFACT_DIR)/coverage.html

.PHONY: docker_build
docker_build:
	@echo "==> builing docker image"
ifdef tag
	docker build --tag $(tag) .
else
	docker build --tag latest .
endif

.PHONY: docker_run
docker_run: docker_build
	@echo "==> builing image and running docker container"
	docker run --rm -d -p 9091:9091 verisart latest
