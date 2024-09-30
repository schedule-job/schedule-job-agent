# Schedule Job Agent

## Run

- Docker

```bash
docker run -it sotaneum/schedule-job-agent:latest -e POSTGRES_SQL_DSN="postgresql://{user}:{pw}@{host}:{port}/{db}?sslmode=disable&search_path={schema}" -e PORT=8080 -e TRUSTED_PROXIES="127.0.0.1,192.168.0.1"
```

- K8S

```bash
git clone https://github.com/schedule-job/schedule-job-agent.git .
cd ./schedule-job-agent/k8s
# set ENV
kubectl apply -f ./deployment.yaml
```

- Local

```bash
git clone https://github.com/schedule-job/schedule-job-agent.git .
cd ./schedule-job-agent
POSTGRES_SQL_DSN="postgresql://{user}:{pw}@{host}:{port}/{db}?sslmode=disable&search_path={schema}" TRUSTED_PROXIES="127.0.0.1,192.168.0.1" PORT=8080 GIN_MODE=release go run .
```

## API

### [POST]/api/v1/request

- body

```json
[
  {
    "id": "UUID",
    "url": "string",
    "method": "string",
    "body": "string",
    "headers": "map[string][]string,"
  },
  ...
]
```

## Build

- Docker

```bash
docker buildx build --push -t sotaneum/schedule-job-agent --platform=linux/amd64,linux/arm64 .
```
