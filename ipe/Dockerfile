# build stage
FROM golang:1.18-alpine AS builder

RUN mkdir /app

WORKDIR /app

RUN apk add git

RUN git clone https://github.com/dimiro1/ipe.git && cd ipe/cmd && go build -o ipe

RUN chmod +x ipe

# stage 2
FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/ipe/cmd/ipe /app
COPY config.yml /app

CMD [ "/app/ipe" ]