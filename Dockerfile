FROM golang:1.13.0-stretch as builder

COPY . /build
WORKDIR /build

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/app .
CMD ["./app"]
