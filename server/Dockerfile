FROM golang:alpine AS build-env

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
FROM alpine
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
