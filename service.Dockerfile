FROM golang:1.21.1 as builder

WORKDIR /app

COPY . .

RUN export GOPROXY=https://goproxy.io \
    && export GO111MODULE=on \
    && go mod tidy \
    && ./scripts/build-service.sh

FROM alpine:latest

WORKDIR /app

RUN set -x \
 && apk add --no-cache ca-certificates bash tini

COPY --from=builder /app/bin/brewing-service /usr/bin/.
COPY --from=builder /app/templates /app/templates

ENV OPENAI_BASE_URL=""
ENV OPENAI_TOKEN=""
ENV WHISPER_ENDPOINT=""
ENV WHISPER_ENDPOINT_SCHEMA=""
ENV OPENAI_PROXY=""
ENV DOCKER_HOST=""
ENV BARK_NOTIFY_DOMAIN=""
ENV BARK_NOTIFY_TOKEN=""
ENV SHARE_DOMAIN=""

EXPOSE 6636

ENTRYPOINT ["tini", "--"]

CMD ["brewing-service"]