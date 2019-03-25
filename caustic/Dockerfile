FROM openjdk:13-jdk-alpine AS build

WORKDIR /usr/src/app
COPY .mvn .mvn
COPY pom.xml mvnw ./
RUN ./mvnw compile
COPY . .
RUN ./mvnw package -Dmaven.test.skip=true

FROM openjdk:8-jre-alpine
WORKDIR /app

COPY --from=build /usr/src/app/target/caustic-*.jar /app/caustic.jar

EXPOSE 8080

CMD ["java", "-jar", "caustic.jar"]
