FROM golang:1.21.1 as builder

WORKDIR /app

COPY . .

RUN export GOPROXY=https://goproxy.io \
    && export GO111MODULE=on \
    && go mod tidy \
    && ./scripts/build-worker.sh

FROM alpine:latest

RUN set -x \
 && apk update \
 && apk upgrade -a \
 && apk add --no-cache \
        ca-certificates \
        curl \
        bash \
        tini \
        ffmpeg \
        python3 \
        py3-mutagen \
 && curl -L https://github.com/yt-dlp/yt-dlp/releases/download/2023.10.13/yt-dlp -o /usr/local/bin/youtube-dl \
 && chmod a+rx /usr/local/bin/youtube-dl \
 && mkdir /downloads \
 && mkdir -p /.cache \
 && chmod 777 /.cache

WORKDIR /downloads

COPY --from=builder /app/bin/brewing-worker /usr/bin/.

ENTRYPOINT ["tini", "--"]

CMD ["brewing-worker", "endpoint", "whisper:9000"]
