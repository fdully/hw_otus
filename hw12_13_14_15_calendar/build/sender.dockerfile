FROM golang:1.14.4

LABEL maintainer="Fdully <fdully@gmail.com>"

# Build Args
ARG APP_NAME=Sender

# Environment Variables
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apt-get update -y && apt-get install -y wget

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz


COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/sender cmd/sender/main.go

# Run the binary.
CMD ["/go/bin/sender"]