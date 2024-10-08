FROM golang:1.23-bullseye AS builder

RUN apt-get update && apt-get install -y upx

WORKDIR /src
COPY . .
RUN make test && make && upx -9 ./bin/kumo

FROM alpine:3.6 AS ca-certificates
RUN apk add -U --no-cache ca-certificates

FROM scratch
USER 1000
WORKDIR /app
COPY --from=ca-certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/bin/kumo /app/kumo

CMD ["/app/kumo"]
