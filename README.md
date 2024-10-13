# Schedule Job Agent

[![Docker Image Build With Push](https://github.com/schedule-job/schedule-job-agent/actions/workflows/docker-image-build-push.yml/badge.svg)](https://github.com/schedule-job/schedule-job-agent/actions/workflows/docker-image-build-push.yml) [![Docker Pulls](https://img.shields.io/docker/pulls/sotaneum/schedule-job-agent?logoColor=fff&logo=docker)](https://hub.docker.com/r/sotaneum/schedule-job-agent) [![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/schedule-job/schedule-job-agent?logo=go&logoColor=fff)](https://go.dev/) [![Docker Image Size (tag)](https://img.shields.io/docker/image-size/sotaneum/schedule-job-agent/latest?logoColor=fff&logo=docker)](https://hub.docker.com/r/sotaneum/schedule-job-agent) [![postgresql](https://img.shields.io/badge/14_or_higher-blue?logo=postgresql&logoColor=fff&label=PostgreSQL)](https://www.postgresql.org/)

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

- request

  - body

    ```json
    [
      {
        "id": "UUID",
        "url": "string",
        "method": "string",
        "body": "string",
        "headers": "map[string][]string,"
      }
    ]
    ```

- response

  - body

    ```json
    {
      "code": 200,
      "data": "ok"
    }
    ```

### [GET]/api/v1/request/:jobId/logs

- request

  - params

    - `jobId` : `uuid`

  - queries
    - `lastId` : `uuid`
    - `limit` : `int`

- response

  - body

    ```json
    {
      "code": 200,
      "data": [
        {
          "id": "364e921f-0f51-4a3b-862e-fa3da033e2f6",
          "jobId": "550e8400-e29b-41d4-a716-446655440000",
          "status": "succeed",
          "requestUrl": "http://localhost:8080/api/v1/request",
          "requestMethod": "POST",
          "responseStatusCode": 400,
          "createdAt": "2024-09-30T00:06:32.518673+09:00"
        },
        {
          "id": "4c3dee96-8241-4535-8008-2f419013518a",
          "jobId": "550e8400-e29b-41d4-a716-446655440000",
          "status": "progress",
          "requestUrl": "http://localhost:8080/api/v1/request",
          "requestMethod": "POST",
          "responseStatusCode": 0,
          "createdAt": "2024-09-30T00:07:45.98302+09:00"
        }
      ]
    }
    ```

### [GET]/api/v1/request/:jobId/log/:id

- request

  - params

    - `jobId` : `uuid`
    - `id` : `uuid`

- response

  - body

    ```json
    {
      "code": 200,
      "data": {
        "id": "0320d779-e258-4a80-a494-e97646a9ab16",
        "jobId": "550e8400-e29b-41d4-a716-446655440000",
        "status": "succeed",
        "requestUrl": "http://localhost:8080/api/v1/request",
        "requestMethod": "POST",
        "responseStatusCode": 400,
        "createdAt": "2024-09-30T00:06:01.889063+09:00",
        "requestHeaders": {},
        "requestBody": "",
        "responseHeaders": {
          "Content-Length": ["64"],
          "Content-Type": ["application/json; charset=utf-8"],
          "Date": ["Sun, 29 Sep 2024 15:06:01 GMT"],
          "Set-Cookie": [
            "go_session_id=abcd; Path=/; Expires=Sun, 06 Oct 2024 15:06:01 GMT; Max-Age=604800; HttpOnly"
          ]
        },
        "responseBody": "{\"code\":400,\"message\":\"잘못된 파라미터 입니다. (EOF)\"}"
      }
    }
    ```

## Build

- Docker

```bash
docker buildx build --push -t sotaneum/schedule-job-agent --platform=linux/amd64,linux/arm64 .
```
