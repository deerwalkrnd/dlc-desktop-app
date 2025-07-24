# Go compiler
GO = go
EXECUTABLE_NAME = dlc-desktop-app.exe
PARSER_EXECUTABLE = parser.exe

# Directories
SRCDIR = src
BINDIR = bin

# Packages
MAIN_PACKAGE = ./cmd/offline/main.go
PARSER_PACKAGE = ./cmd/parser/main.go

# Output paths
BINARY_NAME = $(BINDIR)/$(EXECUTABLE_NAME)
PARSER_BINARY_NAME = $(BINDIR)/$(PARSER_EXECUTABLE)
DATABASE_NAME = ./dlc.sqlite

# Build flags for size optimization
GO_BUILD_FLAGS = -ldflags="-s -w" -trimpath

.PHONY: all setup build clean gui nocgo release remove_db parser

# Default target
all: clean setup build

# Ensure bin directory exists
setup:
	@mkdir -p $(BINDIR)

# Build binary
build:
	$(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Size-optimized binary built: $(BINARY_NAME)"

# Clean build artifacts
clean:
	@rm -rf $(BINDIR)
	@rm -f $(DATABASE_NAME)

# Build GUI binary (Windows)
gui: clean setup
	$(GO) build -ldflags="-s -w -H=windowsgui" -trimpath -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "GUI binary complete: $(BINARY_NAME)"

# Build without CGO
nocgo: clean setup
	CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "CGO-disabled binary built: $(BINARY_NAME)"

# Cross-compiled Windows release
release: clean setup
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Release build complete: $(BINARY_NAME)"

# Remove database file
remove_db:
	@rm -f $(DATABASE_NAME)
	@echo "Removed $(DATABASE_NAME)"

# Build parser binary
parser:
	$(GO) build $(GO_BUILD_FLAGS) -o $(PARSER_BINARY_NAME) $(PARSER_PACKAGE)
	@echo "Parser binary built: $(PARSER_BINARY_NAME)"
