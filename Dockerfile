FROM maven:3.9.8-eclipse-temurin-21

WORKDIR /app

COPY app-backend/ /app/

RUN mvn clean install -U -DskipTests
RUN mv target/*.jar app-backend.jar

EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app-backend.jar"] 
