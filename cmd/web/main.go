package main

import (
	"database/sql"
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

	err := LoadEnvironmentVariables(".env")
	if err != nil {
		log.Fatalln("could not load environment variables", err.Error())
	}

	postgresSqlBuilder := postgre.NewBuilder()
	postgresConnectionURl := os.Getenv("POSTGRE_CONN")
	db, err := BuildDatabasePool(postgresSqlBuilder, postgresConnectionURl)
	if err != nil {
		log.Fatalln("application could not start", err.Error())
	}
	defer db.Close()

	postgreRepo := postgre.NewRepository(db)
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

// LoadEnvironmentVariables loads environment variables from a certain path
func LoadEnvironmentVariables(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	return nil
}

// BuildDatabasePool loads environment variables and does the database connection
func BuildDatabasePool(builder driver.DatabasePoolConnectionBuilder, connectionUrl string) (*sql.DB, error) {

	db, err := driver.CreateDatabaseConnection(builder, connectionUrl)
	if err != nil {
		return nil, err
	}

	err = driver.TestDatabaseConnection(db.GetPool())
	if err != nil {
		return nil, err
	}

	fmt.Println("Postgres database connected")
	return db.GetPool(), nil
}
