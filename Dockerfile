FROM golang:1.25-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates
RUN adduser -D -u 1337 zerohttp

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG GOOS=linux
ARG GOARCH=amd64

ENV GOOS=$GOOS
ENV GOARCH=$GOARCH

RUN mkdir -p /var/cache/certs
RUN chown -R zerohttp:zerohttp /var/cache/certs
RUN chmod -R 700 /var/cache/certs

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o zerohttp ./main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build/zerohttp /zerohttp
COPY --from=builder --chown=1337:1337 /var/cache/certs /var/cache/certs

USER zerohttp

ENTRYPOINT ["/zerohttp"]
