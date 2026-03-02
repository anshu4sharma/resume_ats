# ---------- Builder ----------
FROM golang:1.25.3-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w -buildid=" \
    -o app ./cmd


# ---------- Runtime ----------
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /src/app .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/app"]