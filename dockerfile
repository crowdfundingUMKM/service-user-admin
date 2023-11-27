FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY . /app

RUN go mod download
RUN go mod tidy

# Copy the rest of the application code
COPY . .

RUN go build -o main .

EXPOSE 8081

CMD [ "./main" ]