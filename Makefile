BINARY    := cyberhudd
CTL       := cyberhudctl
CMD       := ./cmd/$(BINARY)
CTLCMD    := ./cmd/$(CTL)
PREFIX    ?= /usr
BINDIR    ?= $(PREFIX)/bin
SVCDIR    ?= /lib/systemd/system

PKGNAME      ?= cyberhud
VERSION      ?= 0.1.0
RELEASE      ?= 1
ARCH         ?= $(GOARCH_PI)
MAINTAINER   ?= CyberHUD Maintainers
DESCRIPTION  ?= CyberHUD daemon and CLI for Raspberry Pi GPIO, STEMMA/QWIIC, and display control
URL          ?= https://github.com/databeast/cyberhud
DEB_MAINTAINER ?= CyberHUD Maintainers <conrad@1211.net>
DEB_VERSION    ?= $(VERSION)-$(RELEASE)

DISTDIR      ?= dist
FONTGEN_DIR   := buildtools/fontgen
FONT_PKG      := ./display/surface/fonts/...
ICONGEN_DIR   := buildtools/gen-icons
MATERIAL_BASE := https://raw.githubusercontent.com/google/material-design-icons/master/variablefont
MATERIAL_TTF  := MaterialSymbolsOutlined[FILL,GRAD,opsz,wght].ttf
MATERIAL_CP   := MaterialSymbolsOutlined[FILL,GRAD,opsz,wght].codepoints

# Default cross-compile target: Raspberry Pi (64-bit)
GOOS_PI   := linux
GOARCH_PI := arm64
SUDO      ?= sudo
PROFILE   ?= waveshare-1.3hat
DAEMON_FLAGS ?=

.PHONY: all build build-pi install test vet generate clean clean-fonts clean-icons \
        run run-headless run-profile list-panels \
        install-service uninstall service-status \
        prepare-debian deb \
        download-icons extract-icons \
        collect-snapshots generate-gallery docs-preview snapshots \
        website-install website-build website-dev website-preview \
        website-test website-test-e2e website-test-a11y website-test-visual website-clean

# --- Primary targets ---

all: vet test build

build: generate
	go build -o $(BINARY) $(CMD)
	go build -o $(CTL) $(CTLCMD)

build-pi: generate
	GOOS=$(GOOS_PI) GOARCH=$(GOARCH_PI) go build -o $(BINARY) $(CMD)
	GOOS=$(GOOS_PI) GOARCH=$(GOARCH_PI) go build -o $(CTL) $(CTLCMD)

# --- Code generation ---

generate: extract-icons
	go generate $(FONT_PKG)

download-icons:
	@mkdir -p $(ICONGEN_DIR)/.cache
	@if [ -f "$(ICONGEN_DIR)/.cache/MaterialSymbolsOutlined.ttf" ] && [ -s "$(ICONGEN_DIR)/.cache/MaterialSymbolsOutlined.ttf" ]; then \
		echo "Icon TTF already cached"; \
	else \
		echo "Downloading Material Symbols TTF..."; \
		curl -gfL -o "$(ICONGEN_DIR)/.cache/MaterialSymbolsOutlined.ttf" \
			"$(MATERIAL_BASE)/$(MATERIAL_TTF)" || \
			{ rm -f "$(ICONGEN_DIR)/.cache/MaterialSymbolsOutlined.ttf"; \
			  echo "ERROR: Failed to download $(MATERIAL_TTF)" >&2; exit 1; }; \
	fi
	@if [ -f "$(ICONGEN_DIR)/.cache/codepoints" ] && [ -s "$(ICONGEN_DIR)/.cache/codepoints" ]; then \
		echo "Icon codepoints already cached"; \
	else \
		echo "Downloading Material Symbols codepoints..."; \
		curl -gfL -o "$(ICONGEN_DIR)/.cache/codepoints" \
			"$(MATERIAL_BASE)/$(MATERIAL_CP)" || \
			{ rm -f "$(ICONGEN_DIR)/.cache/codepoints"; \
			  echo "ERROR: Failed to download $(MATERIAL_CP)" >&2; exit 1; }; \
	fi

extract-icons: download-icons
	@echo "Icon assets ready"

# --- Quality ---

vet: generate
	go vet ./...

test: generate
	go test ./...

# --- Install / Uninstall ---

install: build
	install -Dm755 $(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)
	install -Dm755 $(CTL) $(DESTDIR)$(BINDIR)/$(CTL)

install-service:
	install -Dm644 systemd/$(BINARY).service $(DESTDIR)$(SVCDIR)/$(BINARY).service
	systemctl daemon-reload
	systemctl enable --now $(BINARY).service

uninstall:
	systemctl disable --now $(BINARY).service || true
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY)
	rm -f $(DESTDIR)$(BINDIR)/$(CTL)
	rm -f $(DESTDIR)$(SVCDIR)/$(BINARY).service
	systemctl daemon-reload

# --- Run ---

run: build
	$(SUDO) ./$(BINARY) $(DAEMON_FLAGS)

run-headless: build
	$(SUDO) ./$(BINARY) -nodisplay $(DAEMON_FLAGS)

run-profile: build
	$(SUDO) ./$(BINARY) -panel $(PROFILE) $(DAEMON_FLAGS)

list-panels:
	go run $(CTLCMD) list-panels

service-status:
	$(SUDO) systemctl status $(BINARY).service --no-pager -l
	$(SUDO) journalctl -u $(BINARY).service -n 80 --no-pager

# --- Packaging ---

prepare-debian:
	chmod +x debian/rules debian/cyberhud.config debian/cyberhud.postinst debian/cyberhud.prerm debian/cyberhud.postrm
	@d=$$(date -R); \
	printf "cyberhud (%s) unstable; urgency=medium\n\n  * Automated package build.\n\n -- %s  %s\n" "$(DEB_VERSION)" "$(DEB_MAINTAINER)" "$$d" > debian/changelog

deb: prepare-debian
	TARGET_GOARCH=$(ARCH) dpkg-buildpackage -us -uc -b -a$(ARCH)
	install -d $(DISTDIR)
	cp -f ../cyberhud_$(DEB_VERSION)_*.deb $(DISTDIR)/

# --- Documentation ---

snapshots:
	go test -run PNGSnapshots -timeout 120s ./display/modes/...

collect-snapshots: snapshots
	go run ./buildtools/docsnap/collect

generate-gallery:
	go run ./buildtools/docsnap/gallery

docs-preview: collect-snapshots generate-gallery
	cd ghpages && mkdocs serve

# --- Website (cyberhud.io) ---

website-install:
	cd website && npm ci

website-dev: website-install
	cd website && npm run dev

website-build: website-install
	cd website && npx astro build

website-preview: website-build
	cd website && npm run preview

website-test: website-install
	cd website && npx vitest run

website-test-e2e: website-build
	cd website && npx playwright test

website-test-a11y: website-build
	cd website && npx playwright test tests/e2e/accessibility.spec.ts

website-test-visual: website-build
	cd website && npx playwright test tests/visual/

website-clean:
	rm -rf website/dist website/node_modules

# --- Cleanup ---

clean-fonts:
	rm -f display/surface/fonts/gen_*.go
	rm -rf $(FONTGEN_DIR)/.cache

clean-icons:
	rm -f display/surface/fonts/gen_material_icons*.go
	rm -rf $(ICONGEN_DIR)/.cache

clean: clean-fonts clean-icons
	rm -f $(BINARY)
	rm -f $(CTL)
	rm -rf $(DISTDIR)
