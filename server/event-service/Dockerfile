##
## Build
##
FROM golang:1.16.2 AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /event-service

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
COPY --from=build /event-service /event-service
EXPOSE 8086
ENTRYPOINT ["/event-service"]