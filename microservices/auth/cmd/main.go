package main

import (
	log "backend/logger"
	sessionRepo "backend/microservices/auth/repository/session"
	userRepo "backend/microservices/auth/repository/user"
	protoAuth "backend/microservices/proto/auth"

	"backend/microservices/auth/usecase"
	//"backend/service/csrf/repository"
	"backend/utils"
	"net"

	"google.golang.org/grpc"
)

func main() {

	//Подключение постгрес
	postDB, err := utils.InitPostgresDB()
	if err != nil {
		log.Error(err)
	}
	//Подключение редис
	redisDB, err := utils.InitRedisDB("redis_db_session")
	if err != nil {
		log.Error(err)
	}

	//Попробую 8081
	authListener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Error(err)
	}

	server := grpc.NewServer()

	authUserRepository := userRepo.NewRepository(postDB)
	authSessionRepository := sessionRepo.NewRepository(redisDB)

	authService := usecase.NewService(authUserRepository, authSessionRepository)
	protoAuth.RegisterAuthServer(server,authService)

	log.Info("started auth microservice on 8081")
	err = server.Serve(authListener)
	if err != nil {
		log.Error("serve troubles")
	}

}