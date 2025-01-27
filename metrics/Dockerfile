FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY ./pkg /app/pkg
COPY ./metrics /app/metrics

WORKDIR /app/metrics

RUN CGO_ENABLED=0 GOOS=linux go build -o metrics

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/metrics/metrics /app/metrics

ENTRYPOINT [ "/app/metrics" ]
