# Go compiler
executable_name = dlc-desktop-app
GO = go

# Directories
SRCDIR = src
BINDIR = bin

# Main package
MAIN_PACKAGE = ./cmd/offline/main.go
PARSER_PACKAGE = ./cmd/parser/main.go

# Output binary name
BINARY_NAME = $(BINDIR)/$(executable_name)
DATABASE_NAME = ./dlc.sqlite

PARSER_BINARY_NAME = $(BINDIR)/parser
# Build flags for size optimization
GO_BUILD_FLAGS = -ldflags="-s -w" -trimpath

.PHONY: all
all: clean setup build

.PHONY: setup
setup:
	mkdir -p $(BINDIR)

.PHONY: build
build:
	$(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Size-optimized binary built"

.PHONY: clean
clean:
	rm -rf $(BINDIR)
	rm -f $(DATABASE_NAME)

.PHONY: gui
gui: clean setup
	$(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "GUI binary complete (no GUI mode on Linux)"

.PHONY: nocgo
nocgo: clean setup
	CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "CGO-disabled binary built"

.PHONY: release
release: clean setup
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME).exe $(MAIN_PACKAGE)
	@echo "Windows release binary built"

.PHONY: remove_db
remove_db:
	rm -f $(DATABASE_NAME)

.PHONY: parser
parser:
	$(GO) build $(GO_BUILD_FLAGS) -o $(PARSER_BINARY_NAME) $(PARSER_PACKAGE)
