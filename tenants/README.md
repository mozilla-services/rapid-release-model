# tenants

CLI app for collecting GCPv2 tenant information from our tenant files.

## Overview

The `tenants` CLI app allows users to retrieve deploymentType information from
our tenant files. This should be limited to "argocd", "gha", or "unknown".

## Installation

### Prerequisites

* Go `1.21.3` or newer

### Installation Steps

To install the `tenants` CLI app, run the following command:

```bash
go install github.com/mozilla-services/rapid-release-model/tenants@latest
```

## Usage

### CLI Commands

The main subcommand is `deploymentType`, which retrieves deploymentType information for
our GCPv2 tenants.

```bash
tenants deploymentType -d <path-to-tenants-directory>
```

Example:

```bash
tenants deploymentType -d ../../tenants -f json -o deployment-type.json
```

This command will read tenant information for services from the tenant yaml files in `../../tenants` directory and
parse the `deployment_type` from them.

### CLI Options

The following options are available for `ciplatforms info`:

| Short Option | Long Option     | Description                                      | Default               |
|--------------|-----------------|--------------------------------------------------|-----------------------|
| `-d`         | `--directory`   | Directory containing yaml tenant files           | `tenants`             |
| `-f`         | `--format`      | Output format: csv or json                       | `csv`                 |
| `-o`         | `--output`      | Output file for results                          | `deployment_type.csv` |

## Configuration

You can configure the `tenants` CLI app by passing CLI flags.

### Output File Formats

* **Output File** (`--output`): The output file can be specified in either JSON or CSV format. The `tenants` app will determine the output format based on the file extension (`.json` or `.csv`).