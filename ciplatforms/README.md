# ciplatforms

CLI app for collecting CI platform information from GitHub. 🚀

## Overview

The `ciplatforms` CLI app allows users to retrieve CI configuration details for
GitHub repositories, specifically checking for CI configuration files used by
popular CI platforms such as:

- **CircleCI**
- **GitHub Actions**
- **Taskcluster**

The app processes multiple repositories in batch mode, and it provides options
for configuring request timeouts and request batch sizes.

## Installation

### Prerequisites

* Go `1.21.3` or newer

### Installation Steps

To install the `ciplatforms` CLI app, run the following command:

```bash
go install github.com/mozilla-services/rapid-release-model/ciplatforms@latest
```

## Usage

### Environment Variables

#### GitHub

The following environment variable is required for authentication with the
GitHub GraphQL API.

Please obtain a GitHub API token with the required scopes (e.g., read-only
access to the repositories you intend to query).

Then set the environment variable:

```bash
export CIPLATFORMS_GITHUB_API_TOKEN='[GitHub API token]'
```

### CLI Commands

The main subcommand is `info`, which retrieves CI platform information for
GitHub repositories.

To retrieve CI platform information:

```bash
ciplatforms info --input <input_file> --output <output_file>
```

Example:

```bash
ciplatforms info --input services.csv --output services_ciplatforms.json
```

This command will read repository information for services from `services.csv`,
query GitHub for CI platform configuration details, and output the results to
`services_ciplatforms.json`.

### CLI Options

The following options are available for `ciplatforms info`:

| Short Option | Long Option     | Description                                      | Default                                             |
|--------------|-----------------|--------------------------------------------------|-----------------------------------------------------|
| `-i`         | `--input`       | Input file containing the list of services       | `services.csv`                                      |
| `-o`         | `--output`      | Output file for results                          | `services_ciplatforms.csv`                          |
| `-t`         | `--gh-token`    | GitHub API token for authentication              | `CIPLATFORMS_GITHUB_API_TOKEN` environment variable |
|              | `--timeout`     | Timeout duration for GitHub API requests         | `10s`                                               |
|              | `--batch-size`  | Number of repositories to process per batch      | `50`                                                |


## Configuration

You can configure the `ciplatforms` CLI app by setting environment variables
and/or passing CLI flags.

### Input and Output File Formats

* **Input File** (`--input`): The input file should be a CSV file (e.g., `services.csv`) listing the services/repositories to query. Each entry should include the GitHub owner and repository name.

* **Output File** (`--output`): The output file can be specified in either JSON or CSV format. The `ciplatforms` app will determine the output format based on the file extension (`.json` or `.csv`).

Example `services.csv` input file format:

```csv
service,repository
monitor,mozilla/blurts-server
autoconnect,mozilla-services/autopush-rs
autoendpoint,mozilla-services/autopush-rs
contile,mozilla-services/contile
```

Example JSON output format (`services_ciplatforms.json`):

```json
[
  {
    "service": "monitor",
    "repository": "mozilla/blurts-server",
    "circle_ci": false,
    "gh_actions": true,
    "taskcluster": false,
    "accessible": true,
    "archived": false
  },
  {
    "service": "autoconnect",
    "repository": "mozilla-services/autopush-rs",
    "circle_ci": true,
    "gh_actions": true,
    "taskcluster": false,
    "accessible": true,
    "archived": false
  },
  {
    "service": "autoendpoint",
    "repository": "mozilla-services/autopush-rs",
    "circle_ci": true,
    "gh_actions": true,
    "taskcluster": false,
    "accessible": true,
    "archived": false
  },
  {
    "service": "contile",
    "repository": "mozilla-services/contile",
    "circle_ci": true,
    "gh_actions": false,
    "taskcluster": false,
    "accessible": true,
    "archived": true
  }
]
```

This output provides details about the CI platform configuration status of each
repository, including its accessibility (a repository may be inaccessible if it
does not exist or if the provided authentication token lacks access), whether it
has been archived, and flags indicating the presence of CircleCI, GitHub
Actions, and Taskcluster configuration files.
