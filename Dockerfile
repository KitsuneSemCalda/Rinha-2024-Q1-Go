FROM golang:1.22.0 as build 

WORKDIR /app

ENV TZ=America/SaoPaulo

COPY /app /app/app
COPY go.mod /app/
COPY go.sum /app/ 
COPY main.go /app/ 

RUN go mod tidy && go build -o /server

ENTRYPOINT [ "/server" ]
