FROM golang:1.18 AS builder

RUN apt-get update && apt-get install -y upx

WORKDIR /src
COPY . .
RUN make && upx -9 ./bin/kumo

FROM scratch
USER 1000
WORKDIR /app
COPY --from=builder /src/bin/kumo /app/kumo

CMD ["/app/kumo"]
