FROM golang:1.12-stretch as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a .
RUN md5sum redash-slackbot

FROM chromedp/headless-shell:latest

RUN apt-get update \
  && apt install -y dumb-init \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder /etc/ssl/certs /etc/ssl/certs

COPY --from=builder /app/redash-slackbot /redash-slackbot

ENTRYPOINT ["dumb-init", "--"]

CMD ["/redash-slackbot"]
