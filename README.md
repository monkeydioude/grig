# Grig

Grig is a web-based dashboard and configuration manager. It provides a user-friendly interface to manage and configure various deployment and proxy tools, as well as local system services.

## What does this project do?

Grig acts as a centralized control panel, primarily interfacing with:
- **Josuke**: A webhook-based deployment manager. You can configure hooks, deployments, branches, actions, and commands via Grig's UI.
- **Capybara**: A reverse proxy and service routing configuration tool. You can define proxies, ports, TLS hosts, and service definitions.
- **System Services**: Provides functionality to view and restart local services (like systemd services).

The backend is built with Go, and it serves an HTML interface rendered using [templ](https://github.com/a-h/templ) and styled with [Tailwind CSS](https://tailwindcss.com/).

## Prerequisites

- Go (1.25+)
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

### 5. Production (Docker)

`Dockerfile.prod` builds the CSS, the templ pages and the binary into one image:

```bash
docker build -f Dockerfile.prod -t grig .
docker run -d --name grig --restart unless-stopped -p 6969:6969 \
  -v /etc/grig:/etc/grig \
  -v /etc/systemd/system:/etc/systemd/system \
  -v /run/dbus/system_bus_socket:/run/dbus/system_bus_socket \
  grig
```

It reads `/etc/grig/grig_server.config.yaml`. Mount each file Grig edits at its host path, since it writes symlinks to them. The container runs as root to restart host services over D-Bus, and Grig has no login: only expose it behind an authenticating proxy. See the header of `Dockerfile.prod`.

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

### Applying a saved config

Josuke and Capybara read their config at start. To offer "Save & restart" on a config page, map its file to the systemd unit reading it:

```yaml
restart_on_save:
  /etc/grig/josuke.yaml: josuke.service
  /etc/grig/capybara.yaml: capybara.service
```

Before restarting Josuke, Grig asks its jobs API (`jobs_api.token` in Josuke's config) and refuses while a deployment runs: the config is saved, Josuke is not restarted. Waiting jobs are kept, in Josuke's queue file.

### Creating a service's repository

The "+ Service" page creates a repository on GitHub or Codeberg, with a push webhook to Josuke. A forge is offered when its token is in the environment:

- `GITHUB_TOKEN`: a fine-grained token on "All repositories" with "Administration" and "Webhooks" read and write (one token per owner: your account or an organization). A classic token needs the `repo` scope only.
- `CODEBERG_TOKEN`: a token on "All" repositories with `repository`, `organization` and `user` on read and write (`user` to create in your own account, `organization` in an organization). `CODEBERG_URL` points to another Forgejo instance.

The webhook targets Josuke's hook of type `github`, with its secret, at the public URL set in the config:

```yaml
# where forges reach Josuke: the webhook is this URL + the hook's path
webhook_base_url: https://example.com
```

### Creating a service's database schema

The "+ Service" page also creates a PostgreSQL schema and its user, when `DB_URL` is in the environment: an admin URL (e.g. `postgres://postgres:<password>@127.0.0.1:5432/<db>`) whose user can create roles and schemas. In one transaction, Grig:

- creates the user with the password generated on the page, and fails if the user already exists;
- creates the schema if missing;
- grants `CONNECT`, `TEMPORARY` and `CREATE` on the database (`CREATE` lets the service run `CREATE SCHEMA IF NOT EXISTS` at start);
- grants `USAGE` and `CREATE` on the schema, read/write on its tables and sequences, and the same on future ones.

It then shows the service's `DB_URL` (new user, `search_path` set to the schema). The password is not stored: copy it then.

### Serving under a sub-path

By default Grig serves everything under `/grig`, for a reverse proxy that forwards the prefix without stripping it:

```
:80/grig  ->  :6969/grig
```

Requests to `/` redirect to the base path, and `/healthcheck` stays reachable unprefixed for probes that do not know it. To serve at the root instead, set `base_path: /`.
