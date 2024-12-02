ARG GO_VERSION=1.23

FROM golang:${GO_VERSION} AS builder

WORKDIR /src
COPY go.mod go.sum  ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/raft-api ./cmd
RUN mkdir -p /data && \
    mkdir -p /data/snapshots

FROM gcr.io/distroless/static-debian12

COPY --from=builder /out/raft-api /raft-api
COPY --from=builder /data /data

CMD ["/raft-api"]
