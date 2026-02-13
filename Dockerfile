# Builder stage
FROM golang:1.26.0-alpine3.23 AS builder

WORKDIR /builder/app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -o server ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go

# Final stage
FROM alpine:3.23
WORKDIR /opt/ordersystem

COPY --from=builder /builder/app/server .

EXPOSE 8000
EXPOSE 8080
EXPOSE 50051

CMD ["./server"]
