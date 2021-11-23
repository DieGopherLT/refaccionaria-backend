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

func main() {
	var postgresConnectionURl, port string

	postgresConnectionURl, port = os.Getenv("DATABASE_URL"), os.Getenv("PORT")
	if postgresConnectionURl == "" || port == "" {
		envs, err := LoadEnvironmentVariables(".env")
		if err != nil {
			log.Fatalln("could not load environment variables", err.Error())
		}
		postgresConnectionURl, port = envs["DATABASE_URL"], envs["PORT"]
	}

	postgresSqlBuilder := postgre.NewBuilder()
	db, err := BuildDatabasePool(postgresSqlBuilder, postgresConnectionURl)
	if err != nil {
		log.Fatalln("application could not start", err.Error())
	}
	defer db.Close()

	postgreRepo := postgre.NewRepository(db)
	repo := controller.NewHandlersRepo(postgreRepo)
	controller.SetHandlersRepo(repo)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: Routes(),
	}

	fmt.Println("Server working on port", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln("could not initialize server", err.Error())
		return
	}
}

// LoadEnvironmentVariables loads environment variables from a certain path
func LoadEnvironmentVariables(path string) (map[string]string, error) {
	envs, err := godotenv.Read(path)
	if err != nil {
		return nil, err
	}
	return envs, nil
}

// BuildDatabasePool builds the database pool by using an specific builder
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
