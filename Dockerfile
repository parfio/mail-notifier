# Go build
FROM golang as build

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOPATH=/

WORKDIR /src/mqtt-mail-notifier/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags -s -a -o mqtt-mail-notifier ./cmd/notifier/

# Service definition
FROM alpine

RUN apk add --update libcap tzdata && rm -rf /var/cache/apk/*

COPY --from=build /src/mqtt-mail-notifier/mqtt-mail-notifier /mqtt-mail-notifier

COPY mail-templates /mail-templates

RUN setcap CAP_NET_BIND_SERVICE=+eip /mqtt-mail-notifier

RUN addgroup -g 1000 -S runnergroup && adduser -u 1001 -S apprunner -G runnergroup
USER apprunner

ENTRYPOINT ["/mqtt-mail-notifier"]
