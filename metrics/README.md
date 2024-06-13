# metrics

CLI app to retrieve data for measuring software delivery performance. 📦

## Overview

The `metrics` CLI app currently supports:

* Retrieving data about GitHub Releases
* Retrieving data about GitHub Pull Requests
* Retrieving data about GitHub Deployments
* Retrieving data about Deployments from Grafana Annotations

## Installation

### Prerequisites

* Go `1.21.3` or newer

### Installation Steps

```bash
go install github.com/mozilla-services/rapid-release-model/metrics@latest
```

## Usage

### Environment Variables

The following environment variables are required for authentication with the GitHub GraphQL API and the Grafana REST API.

#### GitHub

For GitHub, please obtain a GitHub API token with the required scopes (e.g. read-only access on public repositories).

Then set the environment variable:

```bash
export RRM_METRICS__GITHUB__TOKEN='[GitHub API token]'
```

#### Grafana

For Grafana, please obtain a Grafana API token with the required access (e.g. read-only access).

Then set the environment variables:

```bash
export RRM_METRICS__GRAFANA__SERVER_URL='[Grafana Server URL]'
export RRM_METRICS__GRAFANA__TOKEN='[Grafana API Token]'
```

### CLI Commands

#### GitHub

For GitHub use:

```bash
metrics github [command]
```

### Grafana

For Grafana use:

```bash
metrics grafana [command]
```

## Configuration

You can configure the `metrics` CLI app by setting environment variables and/or passing CLI flags.

Please note that each subcommand may define additional configuration options via CLI flags. Use `metrics [command] help` to inspect.

### GitHub

For `github prs`, `github releases` and `github deployments`:

| Description              | Environment Variable              | CLI Flags                 |
|--------------------------|-----------------------------------|---------------------------|
| Owner of the GitHub repo | `RRM_METRICS__GITHUB__REPO_OWNER` | `-o, --repo-owner string` |
| Name of the GitHub repo  | `RRM_METRICS__GITHUB__REPO_NAME`  | `-n, --repo-name string`  |

### Grafana

For `grafana deployments`:

| Description                                  | Environment Variable                      | CLI Flags               |
|----------------------------------------------|-------------------------------------------|-------------------------|
| Name of the Grafana app                      | `RRM_METRICS__GRAFANA__ANNOTATIONS__APP`  | `-a, --app-name string` |
| Epoch datetime in milliseconds (e.g. now-6M) | `RRM_METRICS__GRAFANA__ANNOTATIONS__FROM` | `--from string`         |
| Epoch datetime in milliseconds (e.g. now)    | `RRM_METRICS__GRAFANA__ANNOTATIONS__TO`   | `--to string`           |
