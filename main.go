package main

import (
	"GoRestProject/app"
	"GoRestProject/controllers"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) //Middleware'e JWT kimlik Doğrulaması Eklenir
	port := os.Getenv("PORT")         //Enviroment Dosyasından port bilgisi getirilir
	if port == " " {
		port = "8000" //localhost:8000
	}
	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //uygulama localhost:8000 altında istekleri dinlemeye başlar
	if err != nil {
		fmt.Print(err)
	}
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")

	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
}
