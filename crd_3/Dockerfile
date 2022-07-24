FROM golang:1.18.3-alpine3.16

WORKDIR /tmp

COPY ./cmd/operator/operator .
RUN pwd
RUN ls -al /tmp
RUN ls -al /bin

CMD [ "/tmp/operator" ]
