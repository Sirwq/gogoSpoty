# Stage 1 - build
FROM golang:1.26
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gogoSpoty ./cmd/gogoSpoty/

# Stage 2 - final
FROM alpine:latest
WORKDIR /app
COPY --from=0 /app/gogoSpoty .
COPY --from=0 /app/static ./static

EXPOSE 5111 6111
CMD ["./gogoSpoty"]