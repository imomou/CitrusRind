FROM golang:1.12.7-alpine3.10

RUN apk add git zip

WORKDIR /app

COPY . .

RUN go get -d ./...

#https://github.com/serverless/serverless/issues/4710
RUN go build -v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o main .

RUN zip main.zip main
