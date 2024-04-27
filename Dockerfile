FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git ca-certificates
WORKDIR $GOPATH/src/
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/runner
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/healthcheck "github.com/bsquidwrd/TwitchEventSubHandler/healthcheck"


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/runner /go/bin/runner
COPY --from=builder /go/bin/healthcheck /go/bin/healthcheck

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "/go/bin/healthcheck" ]
CMD ["/go/bin/runner"]
