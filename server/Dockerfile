FROM golang:alpine@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS build-env

WORKDIR /app
COPY . /app

# install tesseract with english language data
RUN apk add --no-cache \
    tesseract-ocr \
    tesseract-ocr-dev \
    leptonica-dev \
    g++ \
    tesseract-ocr-data-eng \
    ffmpeg \
    && rm -rf /var/cache/apk/*

# test
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o mqtt ./main.go
FROM alpine@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c
WORKDIR /app
COPY --from=build-env /app/mqtt /app/

RUN apk add --no-cache \
    tesseract-ocr \
    tesseract-ocr-dev \
    leptonica-dev \
    g++ \
    tesseract-ocr-data-eng \
    ffmpeg \
    && rm -rf /var/cache/apk/*


EXPOSE 8080

CMD ["./mqtt"]
