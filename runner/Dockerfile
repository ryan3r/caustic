FROM golang AS build

WORKDIR /go/src
COPY . .
RUN go build main.go

FROM alpine

RUN apk add --no-cache \
	bash \
	curl \
	gcc \
	g++ \
	openjdk8

WORKDIR /app
COPY --from=build /go/src/main .

CMD ["/app/main"]
