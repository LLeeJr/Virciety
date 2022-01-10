# even better: specify versions
FROM golang:1.16.2 AS builder
# Get ready to build the backend executable
RUN mkdir -p /app
# COPY: -from- -to-
COPY . /app/post-service
# set dir for running commands
WORKDIR /app/post-service
# Build backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go get ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /out .

# multi-staging
FROM scratch
# Add backend runnable
COPY --from=builder /out ./
#COPY --from=builder /app/post-service/wait-for-it.sh ./
ENTRYPOINT ["./out"]
# expose ports
EXPOSE 8083