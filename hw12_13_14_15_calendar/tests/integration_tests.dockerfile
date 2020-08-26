FROM golang:1.14.4

LABEL maintainer="Fdully <fdully@gmail.com>"

# Environment Variables
ENV CGO_ENABLED=0

# Set the Current Working Directory inside the container
WORKDIR /app

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

CMD ["go test -v -tags integration ./tests/calendar_test.go"]
