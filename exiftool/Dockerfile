# source: https://github.com/Miljar/exiftool-docker
FROM alpine:3.12
MAINTAINER nitish@jadia.dev

#ENV EXIFTOOL_VERSION=10.20
ENV EXIFTOOL_VERSION=11.99

RUN apk add --no-cache --update perl make openssl
RUN cd /tmp \
    && wget https://exiftool.org/Image-ExifTool-11.99.tar.gz \
    && tar -zxvf Image-ExifTool-${EXIFTOOL_VERSION}.tar.gz \
    && cd Image-ExifTool-${EXIFTOOL_VERSION} \
    && perl Makefile.PL \
    && make test \
    && make install \
    && cd .. \
    && rm -rf Image-ExifTool-${EXIFTOOL_VERSION}


VOLUME /tmp

WORKDIR /tmp

ENTRYPOINT ["exiftool"]
