
# syntax=docker/dockerfile:1

FROM golang:alpine

WORKDIR /app 

COPY go.mod .
COPY go.sum .

RUN go mod download


COPY ./internal ./internal
COPY ./model ./model
COPY ./pkg ./pkg

COPY *.go ./

RUN go build -o /upwork

EXPOSE 8080

CMD ["/upwork"]