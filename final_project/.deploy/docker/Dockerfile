FROM golang:1.23-alpine
ARG APP_NAME
ARG CLI_NAME

WORKDIR /app

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go mod vendor && go build ./cmd/$APP_NAME
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/$CLI_NAME

#FROM alpine:latest

#WORKDIR /app

#COPY --from=builder /app .
