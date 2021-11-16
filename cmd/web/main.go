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

	db, err := RunDatabaseAndEnvVariables()
	if err != nil {
		log.Fatalln("application could not start", err.Error())
		return
	}
	defer db.GetPool().Close()

	postgreRepo := postgre.NewRepository(db.GetPool())
	repo := controller.NewHandlersRepo(postgreRepo)
	controller.SetHandlersRepo(repo)

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

// RunDatabaseAndEnvVariables loads environment variables and does the database connection
func RunDatabaseAndEnvVariables() (driver.DatabasePoolConnectionBuilder, error) {

	// Loading environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("could not load environment variables", err.Error())
		return nil, err
	}

	postgreDbBuilder := postgre.NewBuilder()
	postgreConnectionURl := os.Getenv("POSTGRE_CONN")
	db, err := driver.CreateDatabaseConnection(postgreDbBuilder, postgreConnectionURl)
	if err != nil {
		return nil, err
	}

	err = driver.TestDatabaseConnection(db.GetPool())
	if err != nil {
		return nil, err
	}

	fmt.Println("Postgres database connected")
	return db, nil
}
