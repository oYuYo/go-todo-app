FROM golang:1.23.4-alpine3.21 AS go-app-dev
RUN apk update && apk add git
ENV ROOT=/go/src/app
WORKDIR ${ROOT}
COPY ./app/main.go ${ROOT}
COPY ./app/go.mod ${ROOT}
RUN go mod tidy