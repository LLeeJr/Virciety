##
## Build
##
FROM golang:1.16.2 AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /comment-service

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
COPY --from=build /comment-service /comment-service
EXPOSE 8084
ENTRYPOINT ["/comment-service"]