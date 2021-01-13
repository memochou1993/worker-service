# build stage
FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./client

# final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=builder /app/main ./client/main
COPY --from=builder /app/client/public ./client/public

ENTRYPOINT ./client/main
