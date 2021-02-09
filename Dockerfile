# Choose the golang image as the build base image
FROM golang:1.15-alpine AS build

# Install git for the version string
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Define the directory we should work in
WORKDIR /app

# Download the necessary go modules
COPY go.mod go.sum ./
RUN go mod download

# Build the application
COPY . .
RUN go build \
        -o server \
        -ldflags "\
            -X github.com/x0tf/server/internal/static.ApplicationMode=PROD \
            -X github.com/x0tf/server/internal/static.ApplicationVersion=$(git rev-parse --abbrev-ref HEAD)-$(git log --pretty=format:'%h' -n 1)" \
        ./cmd/server/

# Run the application in an empty alpine environment
FROM alpine:latest
WORKDIR /root
COPY --from=build /app/server .
CMD ["./server"]