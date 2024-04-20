FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/runner
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/healthcheck "github.com/bsquidwrd/TwitchEventSubHandler/healthcheck"


FROM scratch
COPY --from=builder /go/bin/runner /go/bin/runner
COPY --from=builder /go/bin/healthcheck /go/bin/healthcheck

ENTRYPOINT ["/go/bin/runner"]
