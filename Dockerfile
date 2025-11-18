FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

# Build both binaries
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/worker ./cmd/worker

FROM gcr.io/distroless/static-debian12:nonroot

# Copy both binaries
COPY --from=builder /app/bin/server /server
COPY --from=builder /app/bin/worker /worker

EXPOSE 8080

# Default entrypoint for Cloud Run service
ENTRYPOINT [ "/server" ]
