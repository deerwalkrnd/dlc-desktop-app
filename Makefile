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
# Optional command that signs a Windows .exe in place, e.g.
# SIGN=scripts/sign-windows.sh (see that script for the variables it needs).
SIGN     ?=

.PHONY: all build web app tools test lint run demo winres release clean

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
	test -z "$$(gofmt -l cmd internal tools web/embed.go)" || (gofmt -l cmd internal tools web/embed.go; exit 1)
	cd web && npm run check && npm run lint

run: build
	./$(BIN)/$(APP) -library ./DLC

demo:
	./scripts/make-demo-library.sh DLC

# Windows resources (icon, version information, manifest) as .syso files,
# which the Go linker embeds in Windows builds only.
winres:
	go run ./tools/winres -version $(VERSION) -original $(APP).exe \
		-out cmd/dlc-desktop-app/rsrc_windows_amd64.syso
	go run ./tools/winres -version $(VERSION) -original dlc-check.exe \
		-description "DLC Desktop App file-name checker" \
		-out cmd/dlc-check/rsrc_windows_amd64.syso

# Release archives: the binary embeds the web app, so each archive only
# needs the executables, an empty DLC folder and a short README.
release: web winres
	rm -rf $(DIST)
	@for target in windows/amd64 linux/amd64; do \
		os=$${target%/*}; arch=$${target#*/}; ext=; [ $$os = windows ] && ext=.exe; \
		dir=$(DIST)/$(APP)-$(VERSION)-$$os-$$arch; \
		echo "building $$dir"; \
		mkdir -p $$dir/DLC; \
		GOOS=$$os GOARCH=$$arch $(GOBUILD) -o $$dir/$(APP)$$ext ./cmd/dlc-desktop-app || exit 1; \
		GOOS=$$os GOARCH=$$arch $(GOBUILD) -o $$dir/dlc-check$$ext ./cmd/dlc-check || exit 1; \
		if [ -n "$(SIGN)" ] && [ $$os = windows ]; then \
			$(SIGN) $$dir/$(APP)$$ext && $(SIGN) $$dir/dlc-check$$ext || exit 1; \
		fi; \
		cp scripts/release-readme.txt $$dir/README.txt; cp LICENSE $$dir/LICENSE.txt; \
		(cd $(DIST) && zip -qr $(APP)-$(VERSION)-$$os-$$arch.zip $(APP)-$(VERSION)-$$os-$$arch) || exit 1; \
	done
	@ls -lh $(DIST)/*.zip

clean:
	rm -rf $(BIN) $(DIST)
	rm -f cmd/*/rsrc_windows_*.syso
	find web/build -mindepth 1 ! -name .gitkeep -delete 2>/dev/null || true
