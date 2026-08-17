# Documentation Site - Local Development

This directory contains the MkDocs source for the CyberHUD documentation site. The site is built with [MkDocs](https://www.mkdocs.org/) and the [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/) theme, and deployed automatically to GitHub Pages via CI.

## Prerequisites

- Python 3.x

## Setup

Create and activate a virtual environment, then install dependencies:

```sh
# Create a virtual environment
python -m venv .venv

# Activate on Linux/macOS
source .venv/bin/activate

# Activate on Windows
.venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
```

## Local Preview

Start the development server to preview the site locally:

```sh
# From within the ghpages/ directory
mkdocs serve --config-file mkdocs.yml

# Or from the repository root
mkdocs serve --config-file ghpages/mkdocs.yml
```

The site is served at [http://localhost:8000](http://localhost:8000). MkDocs watches for file changes and automatically rebuilds the site, with live reload in the browser so you see updates without manually refreshing.

## Building

To produce a full static build (the same command CI uses):

```sh
mkdocs build --strict --config-file mkdocs.yml
```

Output goes to the `site/` directory. Strict mode catches broken links and missing pages, so if the build succeeds locally it will also succeed in CI.
