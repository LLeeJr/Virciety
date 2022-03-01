##
## Build
##
FROM golang:1.16.2 AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /dm-service

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
COPY --from=build /dm-service /dm-service
EXPOSE 8081
ENTRYPOINT ["/dm-service"]