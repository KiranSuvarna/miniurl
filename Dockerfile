# Stage 1
FROM golang:1.12.8 AS builder
# Creating working directory
RUN mkdir -p  /go/src/bitbucket.org/smartclean
# Copying source code to repository
COPY   .  /go/src/bitbucket.org/smartclean/routines-go
WORKDIR /go/src/bitbucket.org/smartclean/routines-go
# Installing ca certificates
RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates && rm -rf /var/lib/apt/lists/*
ENV GO111MODULE=on
RUN go mod init && go clean
# Creating go binary
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o routines-go
# Stage 2
FROM alpine
RUN apk add --no-cache openssh
# Copy ca certificates from builder
COPY --from=builder  /etc/ssl/certs /etc/ssl/certs
# Copy our static executable and dependencies from builder
COPY --from=builder /go/src/bitbucket.org/smartclean/routines-go  /
COPY --from=builder /go/src/bitbucket.org/smartclean/routines-go/config.yml  /


# Exposing port
EXPOSE 8080
# Run the widget-server  binary.
ENTRYPOINT ["/routines-go"]