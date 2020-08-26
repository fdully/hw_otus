FROM golang:1.14.4 AS builder

LABEL maintainer="Fdully <fdully@gmail.com>"

# Build Args
ARG APP_NAME=goose

# Environment Variables
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY migrations ./migrations

# Build the binary.
RUN go get -u github.com/pressly/goose/cmd/goose

CMD ["/go/bin/goose"]
