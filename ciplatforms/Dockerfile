FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ciplatforms

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/ciplatforms /app/ciplatforms

ENTRYPOINT [ "/app/ciplatforms" ]
