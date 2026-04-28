FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api .

FROM scratch
COPY --from=builder /app/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
