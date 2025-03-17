FROM golang:alpine AS build-env

WORKDIR /app
COPY . /app
RUN go build -o mqtt ./main.go
FROM alpine
WORKDIR /app
COPY --from=build-env /app/mqtt /app/

EXPOSE 8080

CMD ["./mqtt"]
