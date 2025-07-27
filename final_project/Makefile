SHELL := /bin/bash

GO := go
APP_NAME := rate_limiter
CLI_NAME := cli

include .env.docker
export $(shell sed 's/=.*//' .env)
export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

local-tidy:
	go mod tidy
local-vendor:
	go mod vendor
local-build-app:
	go build -o $(APP_NAME) ./cmd/$(APP_NAME)
local-build-cli:
	go build -o $(CLI_NAME) ./cmd/$(CLI_NAME)

local-build: local-tidy local-vendor local-build-app local-build-cli

local-run:
	./$(APP_NAME)
local-clean:
	go clean
	rm $(APP_NAME)
	rm $(CLI_NAME)
local-lint:
	golangci-lint run ./...
local-test:
	go test -race -count 100 ./...
local-all:
	make local-test
	make local-build

local-check:
	./$(CLI_NAME) --conf=.env.docker --cmd=check --login=${login} --pass=${pass} --ip=${ip}
local-reset:
	./$(CLI_NAME) --conf=.env.docker --cmd=reset --login=${login} --ip=${ip}
local-ab:
	./$(CLI_NAME) --conf=.env.docker --cmd=ab --ip=${ip}
local-rb:
	./$(CLI_NAME) --conf=.env.docker --cmd=rb --ip=${ip}
local-aw:
	./$(CLI_NAME) --conf=.env.docker --cmd=aw --ip=${ip}
local-rw:
	./$(CLI_NAME) --conf=.env.docker --cmd=rw --ip=${ip}

docker-check:
	docker exec rate-limiter-ms ./$(CLI_NAME) --conf=.env.docker --cmd=check --login=${login} --pass=${pass} --ip=${ip}
docker-reset:
	docker exec rate-limiter-ms ./$(CLI_NAME) --conf=.env.docker --cmd=reset --login=${login} --ip=${ip}
docker-ab:
	docker exec rate-limiter-ms ./$(CLI_NAME) --conf=.env.docker --cmd=ab --ip=${ip}
docker-rb:
	docker exec rate-limiter-ms ./$(CLI_NAME) --conf=.env.docker --cmd=rb --ip=${ip}
docker-aw:
	docker exec rate-limiter-ms ./$(CLI_NAME) --conf=.env.docker --cmd=aw --ip=${ip}
docker-rw:
	docker exec rate-limiter-ms ./$(CLI_NAME) --conf=.env.docker --cmd=rw --ip=${ip}

docker-test:
	docker exec rate-limiter-ms go test -race -count 100 ./...

compose-up:
	docker compose up -d --build
protoc-install:

protoc:
	protoc --go_out=pb --go-grpc_out=pb api/rate_limiter.proto

.PHONY: local-all local-test local-build


