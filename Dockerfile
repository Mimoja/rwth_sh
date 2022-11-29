FROM golang:alpine
MAINTAINER Mimoja <git@mimoja.de>, Tyalie <git@flowerpot.me>

RUN  apk add --no-cache build-base

RUN mkdir /app
WORKDIR /app
# layer for dependencies 
COPY go.mod go.sum /app/
RUN go mod download

# layer for application code
COPY . /app/
RUN go build -o main go-link-shortener/server

CMD ["/app/main"]
