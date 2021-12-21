.PHONY: build
build:
	make server
	make auth-service
	make event-service
	make user-service

.PHONY: server
server:
	go build -o bin/api/server -v ./cmd/server

.PHONY: auth-service
auth-service:
	go build -o bin/auth-service/auth -v ./cmd/auth

.PHONY: event-service
event-service:	
	go build -o bin/event-service/event -v ./cmd/event

.PHONY: user-service
user-service:
	go build -o bin/user-service/user -v ./cmd/user

.PHONY: cover
cover:
	go test -cover -coverprofile=cover.out -coverpkg=./... ./...
	cat cover.out | fgrep -v "main.go" | fgrep -v "mock.go" | fgrep -v "pb.go" | fgrep -v "response_easyjson.go" > cover1.out
	go tool cover -func=cover1.out

.PHONY:
remove_containers:
	-docker stop $$(docker ps -aq)
	-docker rmi $$(docker images -q)





