FROM openjdk:13-jdk-alpine AS build

WORKDIR /usr/src/app
COPY . .
RUN chmod +x mvnw && \
	sed -i 's/localhost/\$\{MYSQL_HOST\}/g' src/main/resources/application.properties && \
	sed -i 's/Pa$$word1/\$\{MYSQL_ROOT_PASSWORD\}/g' src/main/resources/application.properties && \
	./mvnw package -Dmaven.test.skip=true

FROM openjdk:8-jre-alpine
WORKDIR /app

COPY entrypoint.sh .
COPY --from=build /usr/src/app/target/caustic-*.jar /app/caustic.jar

EXPOSE 8080

CMD ["sh", "entrypoint.sh"]
