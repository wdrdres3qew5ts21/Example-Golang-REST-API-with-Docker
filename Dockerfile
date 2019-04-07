from golang
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN ls
run cd ..
run ls
CMD ["app"]