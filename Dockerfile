FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o main
ENTRYPOINT ["/app/main"]
EXPOSE 8088