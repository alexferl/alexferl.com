FROM abiosoft/caddy:builder as builder

ARG version="0.11.0"
ARG plugins="git"

RUN VERSION=${version} PLUGINS=${plugins} /bin/sh /usr/bin/builder.sh

FROM alpine:3.7
LABEL maintainer "Alexandre Ferland <aferlandqc@gmail.com>"

RUN apk add --no-cache ca-certificates \
    git && \
    update-ca-certificates

LABEL caddy_version="0.11.0"

# install caddy
COPY --from=builder /install/caddy /usr/bin/caddy

# validate install
RUN /usr/bin/caddy -version
RUN /usr/bin/caddy -plugins

COPY Caddyfile /etc/Caddyfile

WORKDIR /srv

EXPOSE 2015

ENTRYPOINT ["/usr/bin/caddy"]
CMD ["--conf", "/etc/Caddyfile", "--log", "stdout"]
