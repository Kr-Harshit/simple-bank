# Build stage
FROM golang:alpine3.18 AS builder
RUN apk add git
RUN apk add curl
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOPROXY=direct go build -o main main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz


# RUN Stage
FROM golang:alpine3.18
WORKDIR /app
COPY app.env .
COPY wait-for.sh .
COPY start.sh .
COPY service/db/migration ./migration
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]

