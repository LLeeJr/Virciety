##
## Build
##
FROM golang:1.16.2 AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /user-service

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
COPY --from=build /user-service /user-service
EXPOSE 8085
ENTRYPOINT ["/user-service"]