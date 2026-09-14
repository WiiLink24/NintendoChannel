FROM golang:alpine AS builder-common

# We assume only git is needed for all dependencies.
# openssl is already built-in.
RUN apk add -U --no-cache git

RUN adduser -D server
USER server
WORKDIR /home/server


# Build file generator CLI
FROM builder-common AS builder-cli

# Cache pulled dependencies if not updated.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy necessary parts of the source into builder's source
COPY cli cli
COPY common common
COPY constants constants
COPY gametdb gametdb
COPY v3 v3
COPY v6 v6

WORKDIR /home/server/cli
# Build to name "app".
RUN go build -o /home/server/app .


# Build Docker UNIX socket interface
FROM builder-common AS builder-docker

# Cache pulled dependencies if not updated.
COPY docker/go.mod .
COPY docker/go.sum .
RUN go mod download

# Copy necessary parts of the source into builder's source
COPY docker/*.go .

# Build to name "app".
RUN go build -o app .


# Runner
FROM alpine:latest

RUN adduser -D server
USER server
WORKDIR /home/server

# Copy executables
COPY --from=builder-cli /home/server/app cli
COPY --from=builder-docker /home/server/app docker

CMD ["./docker"]