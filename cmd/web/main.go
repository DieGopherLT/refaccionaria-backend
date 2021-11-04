package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DieGopherLT/refaccionaria-backend/internal/controller"
	"github.com/DieGopherLT/refaccionaria-backend/internal/driver"
	"github.com/DieGopherLT/refaccionaria-backend/internal/repository/postgre"
	"github.com/joho/godotenv"
)

const PORT = 4000

func main() {

	db, err := Run()
	if err != nil {
		log.Fatalln("application could not start", err.Error())
		return
	}
	defer db.GetPool().Close()

	postgreRepo := postgre.NewRepository(db.GetPool())
	repo := controller.NewRepo(postgreRepo)
	controller.NewHandlers(repo)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: Routes(),
	}

	fmt.Println("Server working on port", PORT)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln("could not initialize server", err.Error())
		return
	}
}

func Run() (driver.SQLDatabase, error) {

	// Loading environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("could not load environment variables", err.Error())
		return nil, err
	}

	// Connection to database
	db, err := driver.ConnectSQL(postgre.NewBuilder(), os.Getenv("POSTGRE_CONN"))
	if err != nil {
		return nil, err
	}
	// Test database connection
	err = driver.TestSQL(db.GetPool())
	if err != nil {
		return nil, err
	}
	fmt.Println("Postgres database connected")

	return db, nil
}
