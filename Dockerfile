FROM golang:1.10-alpine as go

RUN apk update \
	&& apk add git \
	&& rm -rf /var/cache/apk/*

COPY . /go/src/depcharge
WORKDIR /go/src/depcharge

RUN go get ./...
RUN go build
RUN ./depcharge -f -- go get {{get}}

RUN go test .
RUN go build -ldflags="-w -s -X main.version=$VERSION" .

FROM alpine:latest
COPY --from=go /go/src/depcharge/depcharge /bin/depcharge
RUN mkdir /mount
WORKDIR /mount
ENTRYPOINT ["depcharge"]
CMD ["--help"]
