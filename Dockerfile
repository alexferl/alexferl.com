ARG NGINX_VERSION=1.19.3-alpine

FROM nginxinc/nginx-unprivileged:${NGINX_VERSION}

COPY index.html /usr/share/nginx/html

USER root

RUN rm /etc/nginx/conf.d/default.conf
COPY nginx/nginx.conf /etc/nginx/
COPY nginx/alexferl.com.conf /etc/nginx/conf.d/
RUN mkdir -p /etc/nginx/headers.d
COPY nginx/headers.conf /etc/nginx/headers.d/

USER nginx
