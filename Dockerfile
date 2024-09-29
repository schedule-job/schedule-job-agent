FROM golang:alpine as build

WORKDIR /go/
ADD . /go/

RUN apk add --no-cache bash git openssh tzdata ca-certificates
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM scratch as final

EXPOSE 8080/tcp
EXPOSE 8080/udp

ENV TZ=Asia/Seoul

COPY --from=build /go/main .
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/main"]