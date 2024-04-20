FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/runner


FROM scratch AS final
EXPOSE 8080
COPY --from=builder /go/bin/runner /go/bin/runner
ENTRYPOINT ["/go/bin/runner"]
