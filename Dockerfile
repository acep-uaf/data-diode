# syntax=docker/dockerfile:1
# Data Diode CLI in Go via container image.
# https://docs.docker.com/language/golang/build-images/

FROM golang:alpine
WORKDIR /cli
COPY . /cli
RUN CGO_ENABLED=0 GOOS=linux go build -o diode
CMD [ "./diode", "mqtt" ]
