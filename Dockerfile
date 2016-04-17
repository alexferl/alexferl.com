FROM golang:1.6-alpine

RUN echo 'http://dl-4.alpinelinux.org/alpine/v3.3/main' >> /etc/apk/repositories && \
  apk upgrade --update && \
  apk add git

RUN go get github.com/mholt/caddy

WORKDIR /app/

COPY Caddyfile index.html ./

EXPOSE 2015

CMD caddy -log stdout
