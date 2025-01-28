FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY ./pkg /app/pkg
COPY ./ciplatforms /app/ciplatforms

WORKDIR /app/ciplatforms

RUN CGO_ENABLED=0 GOOS=linux go build -o ciplatforms

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/ciplatforms/ciplatforms /app/ciplatforms

ENTRYPOINT [ "/app/ciplatforms" ]
