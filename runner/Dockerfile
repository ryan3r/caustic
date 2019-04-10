FROM golang:1.12-alpine AS build

WORKDIR /go/src
RUN apk add --no-cache git && \
	go get github.com/golang/glog
COPY . .
RUN go build main.go

FROM golang:1.12-alpine

RUN apk add --no-cache \
	bash \
	curl \
	gcc \
	g++ \
	openjdk8 \
	python

WORKDIR /app
ENV PATH="$PATH:/usr/lib/jvm/java-1.8-openjdk/bin/"
COPY --from=build /go/src/main .

CMD ["/app/main"]
