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

COPY cloudreve ./cloudreve

RUN chmod +x ./cloudreve

EXPOSE 5212 443

ENTRYPOINT ["./cloudreve"]

