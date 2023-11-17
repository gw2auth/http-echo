FROM golang:1.19
WORKDIR /go/src/github.com/gw2auth/http-echo/
COPY go.mod ./
COPY main.go ./
RUN go build -o app

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/src/github.com/gw2auth/http-echo/app ./app
CMD ["/root/app"]
