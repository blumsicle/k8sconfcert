FROM golang:1.23 AS build

WORKDIR /app

COPY .git go.mod go.sum Makefile ./
RUN make deps

COPY . .

RUN make build

FROM alpine:latest

COPY --from=build /app/bin/* /app
EXPOSE 8080
ENTRYPOINT ["/app"]
