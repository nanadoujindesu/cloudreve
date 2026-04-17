FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git zip

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Create placeholder assets.zip if not present (frontend can be served externally)
RUN if [ ! -f application/statics/assets.zip ]; then \
      mkdir -p /tmp/assets/build && \
      echo '{}' > /tmp/assets/build/manifest.json && \
      cd /tmp && zip -r /build/application/statics/assets.zip assets/build; \
    fi

RUN CGO_ENABLED=0 go build -o cloudreve .

FROM alpine:latest

WORKDIR /cloudreve

RUN apk update \
    && apk add --no-cache tzdata vips-tools ffmpeg libreoffice font-noto font-noto-cjk libheif libraw-tools\
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

ENV CR_SETTING_DEFAULT_thumb_ffmpeg_enabled=1 \
    CR_SETTING_DEFAULT_thumb_vips_enabled=1 \
    CR_SETTING_DEFAULT_thumb_libreoffice_enabled=1 \
    CR_SETTING_DEFAULT_media_meta_ffprobe=1  \
    CR_SETTING_DEFAULT_thumb_libraw_enabled=1

COPY --from=builder /build/cloudreve ./cloudreve

RUN chmod +x ./cloudreve

EXPOSE 5212 443

ENTRYPOINT ["./cloudreve"]

