FROM golang:alpine as build

RUN apk add ca-certificates 

WORKDIR /opt

COPY . . 

RUN go build -o bin/client cmd/client/main.go

#######################################

FROM alpine:latest

WORKDIR /opt

COPY --from=build /opt/bin/client .

ENTRYPOINT [ "/opt/client"]
