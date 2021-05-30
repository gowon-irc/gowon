FROM golang:alpine as build-env
ADD . /src
RUN cd /src/demo/module1 && go build -o module1

FROM alpine
WORKDIR /app
COPY --from=build-env /src/demo/module1/module1 /app/
ENTRYPOINT ./module1
