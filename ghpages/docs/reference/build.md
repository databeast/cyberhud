# Build and Deployment

This page documents the build toolchain, CI/CD workflows, and external dependencies for CyberHUD.

## Go Version

CyberHUD requires **Go 1.26** or later, as specified by the `go 1.26.0` directive in `go.mod`.

## Makefile Targets

The project Makefile provides the following targets:

| Target | Description |
|--------|-------------|
| `build` | Compiles `cyberhudd` (daemon) and `cyberhudctl` (CLI) binaries for the host platform after running code generation |
| `build-pi` | Cross-compiles both binaries for Raspberry Pi (linux/arm64) after running code generation |
| `generate` | Runs code generation for fonts and icons (downloads Material Symbols assets if not cached, then runs `go generate`) |
| `test` | Runs `go test ./...` across all packages after code generation |
| `vet` | Runs `go vet ./...` across all packages after code generation |
| `install` | Builds then installs both binaries to `$(DESTDIR)/usr/bin/` with mode 755 |
| `deb` | Builds an arm64 Debian package using `dpkg-buildpackage` and copies the `.deb` to `dist/` |
| `clean` | Removes compiled binaries, generated font/icon source files, caches, and the `dist/` directory |
| `run` | Builds then runs the daemon with sudo (supports `DAEMON_FLAGS` for extra arguments) |
| `run-headless` | Builds then runs the daemon with `-nodisplay` flag (no physical display output) |

## CI/CD Workflows

### Debian Packaging (`.github/workflows/deb.yml`)

Builds an arm64 `.deb` package and optionally publishes it as a GitHub release asset.

**Triggers:**

- Push to `main` branch
- Tag push matching `v*`
- Pull requests (all branches)
- Manual dispatch (with optional version override)

**Architecture:** arm64

**Build steps:**

1. Check out repository
2. Set up Go (version from `go.mod`)
3. Derive package version from tag, manual input, or run number
4. Install Debian packaging dependencies (`debhelper`, `dpkg-dev`, `fakeroot`)
5. Run code generation (`make generate`)
6. Run tests (`go test ./...`)
7. Build arm64 Debian package (`make deb`)
8. Upload `.deb` as artifact

**Outputs:**

- Artifact: `cyberhud-deb-arm64` containing the `.deb` file
- On tagged builds (`v*`): the `.deb` is uploaded to the corresponding GitHub Release with auto-generated release notes

---

### MkDocs Site Deployment (`.github/workflows/docs.yml`)

Builds and deploys the MkDocs documentation site to GitHub Pages.

**Triggers:**

- Push to `main` branch affecting `ghpages/`, `display/modes/`, or `tools/docsnap/`
- Pull requests affecting the same paths
- Manual dispatch

**Build steps:**

1. Check out repository
2. Set up Go (version from `go.mod`)
3. Set up Python 3.x and install MkDocs dependencies from `ghpages/requirements.txt`
4. Run code generation (`make generate`)
5. Generate mode PNG snapshots (`make snapshots`)
6. Collect snapshots into docs directory (`make collect-snapshots`)
7. Generate gallery page (`make generate-gallery`)
8. Build site with `mkdocs build --strict`
9. Upload pages artifact (push only)

**Deployment target:** GitHub Pages (via `actions/deploy-pages`), deployed on push to `main` only

---

### Astro Website Deployment (`.github/workflows/website.yml`)

Builds the Astro marketing site (cyberhud.io) and runs Lighthouse CI.

**Triggers:**

- Push to `main` branch affecting `website/`
- Manual dispatch

**Node.js version:** 22.x

**Build steps:**

1. Check out repository
2. Set up Node.js 22.x with npm cache
3. Install dependencies (`npm ci` in `website/`)
4. Build site (`npx astro build` in `website/`)
5. Run Lighthouse CI (desktop and mobile configurations)

## External Dependencies

Direct dependencies declared in `go.mod`:

| Module | Purpose |
|--------|---------|
| `periph.io/x/conn/v3` | Hardware abstraction layer providing interfaces for SPI, I2C, and GPIO communication with display panels and peripherals |
| `periph.io/x/host/v3` | Host platform initialization that detects and registers available hardware drivers for periph.io interfaces |
| `github.com/fogleman/gg` | 2D graphics rendering library used to draw shapes, text, and composited frames for display modes |
| `github.com/go-zeromq/zmq4` | ZeroMQ messaging library enabling the ZMQ display mode to receive and render external data streams |
| `github.com/mdlayher/wifi` | Wi-Fi station information retrieval for the WiFi display mode (signal strength, SSID, link speed) |
| `github.com/tarm/serial` | Serial port communication for the Serial display mode, reading data from UART-connected devices |
| `github.com/fsnotify/fsnotify` | File system event notifications used to watch configuration files for live-reload on changes |
| `pgregory.net/rapid` | Property-based testing framework used for randomized test generation in policy and protocol tests |
| `github.com/srwiley/oksvg` | SVG parsing library that loads Material Symbols icon SVGs for rendering on display surfaces |
| `github.com/srwiley/rasterx` | SVG rasterization engine that renders parsed SVG paths into pixel bitmaps for display output |
| `golang.org/x/image` | Extended image format support and font rendering used by the graphics pipeline |
