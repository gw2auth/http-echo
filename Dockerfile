FROM golang:1.26-trixie AS builder

WORKDIR /build

COPY go.mod ./
RUN go mod download

COPY ../.. ./
RUN CGO_ENABLED=0 go build -o http-echo

FROM scratch
COPY --from=builder /build/http-echo /http-echo
EXPOSE 8080
ENTRYPOINT ["/http-echo"]
