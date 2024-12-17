# Common variables
APP_NAME := go-store
MAIN_DIR := .
EXAMPLES_DIR := ./examples/memory-store
COVERAGE_FILE := cover.out
GO := go

.PHONY: all compile test coverage clean

# Meta goal with everything
all: compile test coverage

# Compile current project
compile:
	@echo "Compiling main module..."
	$(GO) build -v $(MAIN_DIR)/...

	@echo "Compiling example submodules..."
	cd $(EXAMPLES_DIR) && $(GO) build -v ./...

# Execute module tests including examples
test: compile
	@echo "Executing main module tests..."
	$(GO) test -v $(MAIN_DIR)/...

	@echo "Executing examples tests..."
	cd $(EXAMPLES_DIR) && $(GO) test -v ./...

# Generate combined coverage report
coverage: compile
	@echo "Executing main module tests for coverage report..."
	$(GO) test -coverprofile=cover_main.out $(MAIN_DIR)/...

	@echo "Executing examples modules tests for coverage report..."
	cd $(EXAMPLES_DIR) && $(GO) test -coverprofile=cover_examples.out ./...

	@echo "Combining coverage reports..."
	echo "mode: set" > $(COVERAGE_FILE)
	tail -n +2 cover_main.out >> $(COVERAGE_FILE)
	tail -n +2 $(EXAMPLES_DIR)/cover_examples.out >> $(COVERAGE_FILE)

	echo "Generating coverage HTML report..."
	$(GO) tool cover -html=$(COVERAGE_FILE) -o coverage.html

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -f $(COVERAGE_FILE) cover_main.out $(EXAMPLES_DIR)/cover_examples.out coverage.html
	$(GO) clean -cache -testcache -modcache
