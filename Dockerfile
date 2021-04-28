FROM golang:1.16-alpine AS stage0
RUN apk update --no-cache && apk add --no-cache make ca-certificates
RUN mkdir /app
WORKDIR /app
COPY . /app
RUN cd /app && make build

FROM scratch
COPY --from=stage0 /app/git-telegram-bot /
COPY --from=stage0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/git-telegram-bot"]