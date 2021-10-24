package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

const PORT = ":4000"

func main() {
	var err error

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalln("could not load environment variables", err.Error())
		return
	}

	server := http.Server{
		Addr: PORT,
		Handler: Routes(),
	}

	fmt.Println("Server working on port", PORT)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln("could not initialize server", err.Error())
		return
	}
}
