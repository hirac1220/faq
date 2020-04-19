FROM gcr.io/google-appengine/golang
FROM golang:latest
WORKDIR /go/src/app

RUN apt-get update && apt-get install -y \
  && go get -u \
  github.com/kelseyhightower/envconfig \
  github.com/gorilla/mux \
  github.com/rs/cors