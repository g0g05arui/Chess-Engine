# syntax=docker/dockerfile:1

FROM golang:1.23.3-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    if [ -f go.mod ]; then go mod download; fi

COPY . .

RUN go build -o /out/server  && echo "PORT=8080" > /out/.env

FROM gcr.io/distroless/base-debian12 AS runtime

# Non-root security by default
USER nonroot:nonroot

WORKDIR /app

COPY --from=builder /out/server .
COPY --from=builder /out/.env .

EXPOSE 8080

ENTRYPOINT ["/app/server"]
