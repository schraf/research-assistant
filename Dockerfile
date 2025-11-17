FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/research ./cmd/research

FROM golang:1.24 as runtime
COPY --from=builder /app/bin/research /research

EXPOSE 8080

ENTRYPOINT [ "/research" ]
