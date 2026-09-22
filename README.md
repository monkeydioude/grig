# Grig

Grig is a web-based dashboard and configuration manager. It provides a user-friendly interface to manage and configure various deployment and proxy tools, as well as local system services.

## What does this project do?

Grig acts as a centralized control panel, primarily interfacing with:
- **Josuke**: A webhook-based deployment manager. You can configure hooks, deployments, branches, actions, and commands via Grig's UI.
- **Capybara**: A reverse proxy and service routing configuration tool. You can define proxies, ports, TLS hosts, and service definitions.
- **System Services**: Provides functionality to view and restart local services (like systemd services).

The backend is built with Go, and it serves an HTML interface rendered using [templ](https://github.com/a-h/templ) and styled with [Tailwind CSS](https://tailwindcss.com/).

## Prerequisites

- Go (1.23+)
- `make` utility

## How to run it

### 1. Installation & Setup

First, install the necessary dependencies (Tailwind CLI standalone, templ, and gow for live-reloading):

```bash
make install
```

### 2. Development Mode

To run the project locally with live-reloading enabled for Templ templates, Go code, and Tailwind CSS:

```bash
make dev
```
This will start the server and watch for changes. During development, it uses the configuration file located at `./test/test.config.yml`.

### 3. Compiling and Running

To compile the Tailwind CSS and Templ templates without starting the watcher:

```bash
make compile
```

You can then run the built server directly:

```bash
go run cmd/grig-server/*.go -c grig_server.config.yaml
```
*(You can pass the `-c` flag to point to a custom configuration file, otherwise it defaults to `grig_server.config.yaml`)*

You can also specify a custom port using the `PORT` environment variable:
```bash
PORT=3000 go run cmd/grig-server/*.go
```

### 4. Running via Docker

You can also run Grig using Docker Compose:

```bash
make dev-linux
```

## Configuration

Grig's main configuration file (`grig_server.config.yaml`) is YAML, and points to the Capybara and Josuke config files. If those files do not exist, Grig creates them on startup.

Capybara and Josuke config files are read and written in the format matching their extension: `.yaml`/`.yml` as YAML, anything else as JSON. You can mix the two.

```yaml
# Sub-path the app is served under. Defaults to /grig.
# Set to "/" to serve at the root instead.
base_path: /grig

josuke_config_path: /etc/grig/josuke.config.json

# One entry per Capybara config. Each gets its own page and menu
# entry, named after the file: garden.capybara.yaml shows as "garden".
capybara_config_paths:
    - /etc/grig/capybara.yaml
    - /etc/grig/garden.capybara.yaml

services_paths:
    - /etc/systemd/system/josuke.service
```

### Serving under a sub-path

By default Grig serves everything under `/grig`, for a reverse proxy that forwards the prefix without stripping it:

```
:80/grig  ->  :6969/grig
```

Requests to `/` redirect to the base path, and `/healthcheck` stays reachable unprefixed for probes that do not know it. To serve at the root instead, set `base_path: /`.
