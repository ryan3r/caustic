# syntax = docker/dockerfile:1.0-experimental
FROM openjdk:13-jdk-alpine AS build

WORKDIR /usr/src/app
COPY . .
RUN --mount=type=cache,id=mvn,target=/usr/src/app/target \
	chmod +x mvnw && \
	sed -i 's/localhost/\$\{MYSQL_HOST\}/g' src/main/resources/application.properties && \
	sed -i 's/Pa$$word1/\$\{MYSQL_ROOT_PASSWORD\}/g' src/main/resources/application.properties && \
	sed -i 's/root/\$\{MYSQL_USER\}/g' src/main/resources/application.properties && \
	./mvnw -B package -Dmaven.repo.local=target/.m2/repository -DskipTests && \
	cp target/caustic-*.jar .

FROM openjdk:8-jre-alpine
WORKDIR /app

COPY entrypoint.sh .
COPY --from=build /usr/src/app/caustic-*.jar /app/caustic.jar

EXPOSE 8080

CMD ["sh", "entrypoint.sh"]
