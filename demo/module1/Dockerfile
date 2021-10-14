FROM golang:alpine as build-env
COPY . /src
WORKDIR /src
RUN go build -o module1

FROM alpine:3.14.2
WORKDIR /app
COPY --from=build-env /src/module1 /app/
ENTRYPOINT ./module1
