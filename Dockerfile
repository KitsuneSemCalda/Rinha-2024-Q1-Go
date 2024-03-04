FROM golang:1.22.0-alpine as build 

WORKDIR /app

ENV TZ=America/SaoPaulo
ENV CGO_ENABLED=0

COPY /app app
COPY go.mod .
COPY go.sum .
COPY main.go .

RUN go mod tidy && go build -o /server

ENTRYPOINT [ "/server" ]
