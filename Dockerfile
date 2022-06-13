FROM golang:1.18 AS builder

WORKDIR /opt/app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.16
WORKDIR /opt/app
COPY --from=builder /opt/app/server .
CMD ["./server"]
