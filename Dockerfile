# syntax=docker/dockerfile:1

from golang:1.19.2-bullseye

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY internal ./internal
RUN go mod tidy
RUN go mod download

RUN go build -o /hit-tracker

EXPOSE 9090

CMD [ "/hit-tracker" ]

