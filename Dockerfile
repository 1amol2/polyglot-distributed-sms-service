
# Stage 1: Build the Spring Boot application
FROM gradle:9.7.1-jdk21 AS builder

WORKDIR /workspace

# Copy project files into the build container
COPY --chown=gradle:gradle . .

# Compile the application and create the executable Spring Boot JAR
RUN gradle --no-daemon clean bootJar \
    && find build/libs -maxdepth 1 -name '*.jar' ! -name '*-plain.jar' \
       -exec cp {} /tmp/app.jar \; \
    && test -f /tmp/app.jar

# Stage 2: Run the application without Gradle
FROM eclipse-temurin:21-jre

WORKDIR /app

COPY --from=builder /tmp/app.jar /app/app.jar

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "/app/app.jar"]
