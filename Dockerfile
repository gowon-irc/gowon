FROM golang:alpine as build-env
ADD . /src
RUN cd /src && go build -o gowon

FROM alpine
WORKDIR /app
COPY --from=build-env /src/gowon /app/
ENTRYPOINT ./gowon
