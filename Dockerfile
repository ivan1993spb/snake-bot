
# See Makefile
ARG IMAGE_GOLANG=golang:1.21.5-alpine3.19
ARG IMAGE_ALPINE=alpine:3.19

FROM $IMAGE_ALPINE AS helper

RUN adduser -u 10001 -h /dev/null -H -D -s /sbin/nologin snake

RUN sed -i '/^snake/!d' /etc/passwd

FROM $IMAGE_GOLANG AS builder

ARG VERSION=unknown
ARG BUILD=unknown

WORKDIR /go/src/snake-bot

COPY . .

ENV CGO_ENABLED=0

RUN go build \
    -ldflags "-s -w -X main.Version=${VERSION} -X main.Build=${BUILD::7}" \
    -v -x -o /snake-bot ./cmd/snake-bot

FROM scratch

COPY --from=helper /etc/passwd /etc/passwd

USER snake

COPY --from=builder /snake-bot /usr/local/bin/snake-bot

ENTRYPOINT ["snake-bot"]
