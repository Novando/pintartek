FROM golang:1.22.2-alpine3.19 as src-img

WORKDIR /app

COPY . /app
RUN apk add --no-cache make
RUN make swagger
RUN go build -o ./bin/passvault-service ./cmd/passvault-service/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=src-img /app/bin/passvault-service /app/bin/passvault-service
COPY ./config/config.local.json /app/config/config.local.json
RUN chmod -R 777 /app
RUN ls -l /app

CMD ["/app/bin/passvault-service"]
