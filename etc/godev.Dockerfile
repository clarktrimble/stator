FROM golang:1.21.6-alpine3.19

LABEL description="Golang development"

RUN apk add --no-cache \
    curl \
    git \
    jq \
    make \
    openssh-client

WORKDIR /project

ARG userid
RUN addgroup -S -g ${userid} devo
RUN adduser -S -u ${userid} -g devo devo
USER devo

RUN go install github.com/matryer/moq@v0.3.3
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b $(go env GOPATH)/bin v1.55.2

CMD sh
