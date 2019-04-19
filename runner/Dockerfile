FROM golang:1.12-alpine AS build

WORKDIR /go/src
RUN apk add --no-cache git && \
	go get github.com/docker/docker/client && \
	go get github.com/go-sql-driver/mysql && \
	go get github.com/ahmetalpbalkan/dlog
COPY . .
RUN go build *.go

FROM golang:1.12-alpine

WORKDIR /app
ENV PATH="$PATH:/usr/lib/jvm/java-1.8-openjdk/bin/"
COPY --from=build /go/src/db main

CMD ["/app/main"]
