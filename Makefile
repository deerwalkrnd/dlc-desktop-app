# DLC Desktop App
#
#   make            build the web app and both binaries into bin/
#   make test       run Go and frontend tests
#   make lint       vet, type-check and lint everything
#   make run        build and start the app
#   make demo       create a demo DLC/ library with playable sample videos
#   make release    cross-compile release archives into dist/
#   make clean      remove build output

APP      := dlc-desktop-app
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS  := -s -w -X main.version=$(VERSION)
GOBUILD  := CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)"
BIN      := bin
DIST     := dist

.PHONY: all build web app tools test lint run demo release clean

all: build

build: web app tools

web/node_modules: web/package-lock.json
	cd web && npm ci
	@touch $@

web: web/node_modules
	cd web && npm run build

app:
	$(GOBUILD) -o $(BIN)/$(APP) ./cmd/dlc-desktop-app

tools:
	$(GOBUILD) -o $(BIN)/dlc-check ./cmd/dlc-check

test: web/node_modules
	go test -race ./...
	cd web && npm test

lint: web/node_modules
	go vet ./...
	test -z "$$(gofmt -l cmd internal web/embed.go)" || (gofmt -l cmd internal web/embed.go; exit 1)
	cd web && npm run check && npm run lint

run: build
	./$(BIN)/$(APP) -library ./DLC

demo:
	./scripts/make-demo-library.sh DLC

# Release archives: the binary embeds the web app, so each archive only
# needs the executables, an empty DLC folder and a short README.
release: web
	rm -rf $(DIST)
	@for target in windows/amd64 linux/amd64; do \
		os=$${target%/*}; arch=$${target#*/}; ext=; [ $$os = windows ] && ext=.exe; \
		dir=$(DIST)/$(APP)-$(VERSION)-$$os-$$arch; \
		echo "building $$dir"; \
		mkdir -p $$dir/DLC; \
		GOOS=$$os GOARCH=$$arch $(GOBUILD) -o $$dir/$(APP)$$ext ./cmd/dlc-desktop-app || exit 1; \
		GOOS=$$os GOARCH=$$arch $(GOBUILD) -o $$dir/dlc-check$$ext ./cmd/dlc-check || exit 1; \
		cp scripts/release-readme.txt $$dir/README.txt; \
		(cd $(DIST) && zip -qr $(APP)-$(VERSION)-$$os-$$arch.zip $(APP)-$(VERSION)-$$os-$$arch) || exit 1; \
	done
	@ls -lh $(DIST)/*.zip

clean:
	rm -rf $(BIN) $(DIST)
	find web/build -mindepth 1 ! -name .gitkeep -delete 2>/dev/null || true
