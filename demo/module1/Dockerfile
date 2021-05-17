FROM golang:alpine as build-env
ADD . /src
RUN cd /src && go build -o module1

FROM alpine
WORKDIR /app
COPY --from=build-env /src/module1 /app/
ENTRYPOINT ./module1
