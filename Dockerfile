from golang:1.8-alpine
WORKDIR /go/src/app
run apk update && apk add git
COPY app .
CMD ["./app"]
