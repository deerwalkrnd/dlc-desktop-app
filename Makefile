
# Go compiler
executable_name = "dlc-desktop-app.exe"
GO = go

# Directories
SRCDIR = src
BINDIR = bin

# Main package - update with your entry point
MAIN_PACKAGE = ./cmd/offline/main.go
PARSER_PACKAGE = ./cmd/parser/main.go

# Output binary name
BINARY_NAME = $(BINDIR)/$(executable_name)
DATABASE_NAME = "./dlc.sqlite"

PARSER_BINARY_NAME = $(BINDIR)/parser.exe
# Build flags for size optimization
# -ldflags="-s -w": 
#   -s: disable symbol table
#   -w: disable DWARF generation
# -trimpath: remove file system paths from binary
GO_BUILD_FLAGS = -ldflags="-s -w" -trimpath

.PHONY: all
all: clean setup build 

.PHONY: setup
setup:
	@if not exist $(BINDIR) mkdir $(BINDIR)

.PHONY: build
build:
	$(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo Size-optimized binary built

.PHONY: clean
clean:
	@if exist $(BINDIR) rmdir /S /Q $(BINDIR)
	@if exist $(DATABASE_NAME) del $(DATABASE_NAME)

.PHONY: gui
gui: clean setup
	$(GO) build -ldflags="-s -w -H=windowsgui" -trimpath -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo GUI binary complete.

.PHONY: nocgo
nocgo: clean setup
	set CGO_ENABLED=0 && $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo CGO-disabled binary built.

.PHONY: release
release: clean setup
	set GOOS=windows && set GOARCH=amd64 && set CGO_ENABLED=0 && $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo Release build complete.

.PHONY: remove_db
remove_db:
	rm ./dlc.sqlite

.PHONY: parser
parser:
	$(GO) build $(GO_BUILD_FLAGS) -o $(PARSER_BINARY_NAME) $(PARSER_PACKAGE)

	