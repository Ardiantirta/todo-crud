FROM alpine:latest

RUN mkdir /app

ADD . /app/

WORKDIR /app

ENTRYPOINT  ["/app/main"]