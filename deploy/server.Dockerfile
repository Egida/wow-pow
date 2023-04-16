FROM golang:alpine as build

RUN apk add ca-certificates 

WORKDIR /opt

COPY . . 

RUN go build -o bin/server cmd/server/main.go

#######################################

FROM alpine:latest

WORKDIR /opt

COPY --from=build /opt/bin/server .
COPY --from=build /opt/data ./data

ENTRYPOINT [ "/opt/server"]
