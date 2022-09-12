# build-stage
FROM alpine:3 as build-stage

# use docker build --build-arg VERSION=r2022.3.1.x .
ARG VERSION=r2022.3.1.x
ARG SERVER_VERSION=${VERSION}
ARG CORTEZA_SERVER_PATH=./build/corteza-server-${SERVER_VERSION}.tar.gz

RUN mkdir /tmp/server

ADD $CORTEZA_SERVER_PATH /tmp/server

RUN apk update && apk add --no-cache file

RUN file "/tmp/server/$(basename $CORTEZA_SERVER_PATH)" | grep -q 'gzip' && \
    tar zxvf "/tmp/server/$(basename $CORTEZA_SERVER_PATH)" -C / || \
    cp -a "/tmp/server" /

RUN mv /tmp/server/corteza-server /corteza

WORKDIR /corteza

# deploy-stage
FROM ubuntu:20.04

RUN apt-get -y update \
 && apt-get -y install \
    ca-certificates \
    curl \
 && rm -rf /var/lib/apt/lists/*

ENV STORAGE_PATH "/data"
ENV CORREDOR_ADDR "corredor:80"
ENV HTTP_ADDR "0.0.0.0:80"
ENV HTTP_WEBAPP_ENABLED "false"
ENV PATH "/corteza/bin:${PATH}"

WORKDIR /corteza

VOLUME /data

COPY --from=build-stage /corteza ./

HEALTHCHECK --interval=30s --start-period=1m --timeout=30s --retries=3 \
    CMD curl --silent --fail --fail-early http://127.0.0.1:80/healthcheck || exit 1

EXPOSE 80

ENTRYPOINT ["./bin/corteza-server"]

CMD ["serve-api"]
