FROM golang:1.16-alpine

ENV GO111Module=on

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /post-service

EXPOSE 8083

CMD [ "/post-service" ]