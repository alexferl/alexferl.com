FROM golang:1.6-alpine

RUN echo 'http://dl-4.alpinelinux.org/alpine/v3.4/main' >> /etc/apk/repositories && \
  apk upgrade --update && \
  apk add git curl

RUN curl --silent --show-error --fail --location \
      --header "Accept: application/tar+gzip, application/x-gzip, application/octet-stream" -o - \
      "https://caddyserver.com/download/build?os=linux&arch=amd64&features=git" \
    | tar --no-same-owner -C /usr/bin/ -xz caddy \
 && chmod 0755 /usr/bin/caddy \
 && /usr/bin/caddy -version

WORKDIR /app/

COPY Caddyfile index.html ./

EXPOSE 2015

ENTRYPOINT ["caddy"]
CMD ["-log", "stdout"]
