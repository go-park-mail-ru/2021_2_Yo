package main

import (
	deliveryAuth "backend/auth/delivery/http"
	localStorageAuth "backend/auth/repository/localstorage"
	useCaseAuth "backend/auth/usecase"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)



func main() {
	fmt.Println("Hello, World!")

	r := mux.NewRouter()

	repo := localStorageAuth.NewRepositoryUserLocalStorage()
	useCase := useCaseAuth.NewUseCaseAuth(repo)
	handler := deliveryAuth.NewHandlerAuth(useCase)

	r.HandleFunc("/signup", handler.SignUp)
	r.HandleFunc("/signin", handler.SignIn)

	log.Println("start serving :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}

}
