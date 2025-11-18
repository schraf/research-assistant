FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/research ./cmd/research

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/bin/research /research

EXPOSE 8080

ENTRYPOINT [ "/research" ]
