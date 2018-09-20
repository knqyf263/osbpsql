From golang:1.10-alpine3.7 as builder
RUN apk add --update --no-cache curl git gcc libc-dev build-base
ADD Gopkg.toml Gopkg.lock $GOPATH/src/github.com/knqyf263/osbpsql/
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR $GOPATH/src/github.com/knqyf263/osbpsql/
RUN dep ensure -vendor-only
ADD . $GOPATH/src/github.com/knqyf263/osbpsql/
RUN go build -o osbpsql

From alpine:latest
EXPOSE 8080
COPY --from=builder /go/src/github.com/knqyf263/osbpsql/osbpsql /app/osbpsql
COPY config.toml /app/
WORKDIR /app
CMD ["/app/osbpsql"]
