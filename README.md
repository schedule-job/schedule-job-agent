# Schedule Job Agent

[![Docker Image Build With Push](https://github.com/schedule-job/schedule-job-agent/actions/workflows/docker-image-build-push.yml/badge.svg)](https://github.com/schedule-job/schedule-job-agent/actions/workflows/docker-image-build-push.yml) [![Docker Pulls](https://img.shields.io/docker/pulls/sotaneum/schedule-job-agent?logoColor=fff&logo=docker)](https://hub.docker.com/r/sotaneum/schedule-job-agent) [![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/schedule-job/schedule-job-agent?logo=go&logoColor=fff)](https://go.dev/) [![Docker Image Size (tag)](https://img.shields.io/docker/image-size/sotaneum/schedule-job-agent/latest?logoColor=fff&logo=docker)](https://hub.docker.com/r/sotaneum/schedule-job-agent) [![postgresql](https://img.shields.io/badge/14_or_higher-blue?logo=postgresql&logoColor=fff&label=PostgreSQL
)](https://www.postgresql.org/)



- Immediately sends an API request based on the information received from the API request.
- The requested response results are saved in the DB.

- API 요청 받은 정보를 바탕으로 즉시 API 요청을 보냅니다.
- 요청한 응답 결과를 DB에 저장합니다.

## Require

- PostgreSQL 14 or higher
- Go 1.20 or higher

## Prepared

- [Create Table in PostgreSQL](./sql/agent.sql)

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
]
```

## Build

- Docker

```bash
docker buildx build --push -t sotaneum/schedule-job-agent --platform=linux/amd64,linux/arm64 .
```
