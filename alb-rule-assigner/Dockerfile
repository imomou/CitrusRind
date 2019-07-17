FROM golang:1.12.7-alpine3.10

RUN apk add git

RUN apk add zip

WORKDIR /app

COPY . .

RUN go get -d ./...

RUN go build -o main .

RUN zip main.zip main
