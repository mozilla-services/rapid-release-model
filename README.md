# rapid-release-model

Rapid Release Model for Mozilla Services 🚀

## Docker

### Build the `metrics` Docker image

To build the `metrics` Docker image, run the following in the repository root
directory.

```bash
docker build -t metrics -f metrics/Dockerfile .
```

### Run the `metrics` Docker Container

Run the `github` subcommand inside the Docker container with the required
environment variable:


```bash
docker run --rm \
  -e RRM_METRICS__GITHUB__TOKEN \
  metrics:latest \
  github
```

Run the `grafana` subcommand inside the Docker container with the required
environment variables:

```bash
docker run --rm \
  -e RRM_METRICS__GRAFANA__SERVER_URL \
  -e RRM_METRICS__GRAFANA__TOKEN \
  metrics:latest \
  grafana
```

### Build the `ciplatforms` Docker image

To build the `ciplatforms` Docker image, run the following in the repository
root directory.

```bash
docker build -t ciplatforms -f ciplatforms/Dockerfile .
```

### Run the `ciplatforms` Docker Container

To run the `info` subcommand inside the Docker container, mount the directory
containing the input and output files as a volume and pass the required
environment variable:

```bash
docker run --rm \
  -e CIPLATFORMS_GITHUB_API_TOKEN \
  -v $(pwd)/data:/data \
  ciplatforms:latest \
  info --input /data/services.csv --output /data/services_out.json
```

### Notes

The examples assume the required environment variables are already set in your
host environment.

If the environment variables are **not set on the host**, you must explicitly
pass their values when running the container:

```bash
docker run --rm \
  -e RRM_METRICS__GITHUB__TOKEN=<your_github_token> \
  metrics:latest \
  github
```
