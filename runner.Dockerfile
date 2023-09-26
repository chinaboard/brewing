FROM golang as builder

WORKDIR /app

COPY . .

RUN export GOPROXY=https://goproxy.io \
    && export GO111MODULE=on \
    && go mod tidy \
    && ./scripts/build-runner.sh

FROM alpine:latest

RUN set -x \
 && apk add --no-cache ca-certificates curl bash tini ffmpeg python3 \
 && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/youtube-dl \
 && chmod a+rx /usr/local/bin/youtube-dl \
 && mkdir /downloads \
 && mkdir -p /.cache \
 && chmod 777 /.cache

WORKDIR /downloads

COPY --from=builder /app/bin/brewing-runner /usr/bin/.

ENTRYPOINT ["tini", "--"]

CMD ["brewing-runner","endpoint", "whisper:9000"]