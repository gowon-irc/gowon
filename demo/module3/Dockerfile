FROM golang:alpine as build-env
ADD . /src
RUN cd /src/demo/module3 && go build -o module3

FROM alpine
WORKDIR /app
COPY --from=build-env /src/demo/module3/module3 /app/
ENTRYPOINT ./module3
