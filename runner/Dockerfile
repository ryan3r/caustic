FROM alpine

RUN apk add --no-cache \
	bash \
	curl \
	gcc \
	g++ \
	openjdk8

WORKDIR /app
COPY . .

CMD ["/bin/bash", "./runner.sh"]
