FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o SongsLibrary ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /root/

COPY --from=build /app/SongsLibrary .

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

EXPOSE ${APP_PORT}

CMD ["bash", "/wait-for-it.sh", "db:5432", "--", "./SongsLibrary"]