FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/adewaleafolabi/listing

COPY Gopkg.lock Gopkg.toml ./
COPY vendor vendor
COPY db db
COPY model model
COPY listing-service listing-service

RUN go install ./...

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .