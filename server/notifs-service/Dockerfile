##
## Build
##
FROM golang:1.16.2 AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /notifs-service

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
COPY --from=build /notifs-service /notifs-service
EXPOSE 8082
ENTRYPOINT ["/notifs-service"]