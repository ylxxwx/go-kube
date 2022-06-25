FROM golang:1.18.3-alpine3.16 AS builder

RUN apk update && \
  apk add --no-cache --update make bash git ca-certificates && \
  update-ca-certificates

WORKDIR /go/src/go-kube

COPY . . 


RUN go build -o nmw ./nswatch-client-go/*.go

FROM alpine:3.16.0

COPY --from=builder /go/src/go-kube/nmw /nmw

CMD [ "/nmw" ]
