FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/nuricms
RUN CGO_ENABLED=1 GOOS=linux go build -o nuricms

FROM alpine:latest
RUN apk add --no-cache sqlite-libs
WORKDIR /app
COPY --from=builder /app/cmd/nuricms/nuricms .
ENV GIN_MODE=release
EXPOSE 8080
CMD ["./nuricms"]
