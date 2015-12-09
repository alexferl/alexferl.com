FROM golang:latest

RUN go get github.com/mholt/caddy

RUN apt-get -qq update && apt-get -qq install --no-install-recommends -y nodejs \
npm

RUN ln -s /usr/bin/nodejs /usr/bin/node

RUN npm install --silent bower -g

WORKDIR /app/

COPY bower.json ./

RUN bower --silent install --allow-root

COPY Caddyfile index.html ./

CMD caddy -log stdout -agree -email aferlandqc@gmail.com

