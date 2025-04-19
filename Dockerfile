FROM golang:1.24 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download

RUN GOOS=linux go build

FROM debian:latest AS build-release-stage

WORKDIR /app

COPY --from=build-stage app/webserver ./webserver

EXPOSE 8080

CMD ["./webserver"]