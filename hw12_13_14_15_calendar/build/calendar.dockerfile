FROM golang:1.14.4 AS builder

LABEL maintainer="Fdully <fdully@gmail.com>"

# Build Args
ARG APP_NAME=Calendar

# Environment Variables
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY .. .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/calendar cmd/calendar/main.go


EXPOSE 8080


RUN useradd -u 10001 scratchuser
USER scratchuser

############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/calendar /go/bin/calendar
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# User
COPY --from=0 /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Moscow
USER scratchuser
# Run the binary.
ENTRYPOINT ["/go/bin/calendar"]