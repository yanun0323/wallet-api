##
## Build
##
FROM golang:buster AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /wallet-api

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /wallet-api /wallet-api

EXPOSE 8080

USER root

ENTRYPOINT ["/wallet-api"]