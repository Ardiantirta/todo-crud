FROM alpine:latest

RUN mkdir /app

ADD . /app/

WORKDIR /app

CMD ["/app/main"]