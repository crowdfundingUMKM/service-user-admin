#syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
COPY . /app

RUN go mod download
RUN go mod tidy

RUN go build -o main

EXPOSE 8081

CMD [ "./main" ]