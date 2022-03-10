##
## Build
##
FROM golang:1.16.2 AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /post-service

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
COPY --from=build /post-service /post-service
EXPOSE 8083
ENTRYPOINT ["/post-service"]