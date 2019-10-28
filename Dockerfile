FROM alpine:latest

RUN mkdir /app

ADD . /app/

WORKDIR /app

COPY . .

CMD ["/app/main"]