# syntax=docker/dockerfile:1

FROM golang:1.19-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./ /app/
RUN go build -o /d2lib-docker

EXPOSE 8090

CMD [ "/d2lib-docker" ]
