package main

import (
	"cvital/db"
	"cvital/domain/users"
	"cvital/router"
	"log"
	"net/http"
)

func main() {
	newDb, err := db.NewConnection(db.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "cvital",
		DbName:   "cvital",
		Password: "cvital",
		SslMode:  "disable",
	})
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	err = db.RunMigrations(newDb.DB)
	if err != nil {
		log.Fatalf("DB migrations failed: %v", err)
	}

	usersUseCase := users.NewUseCase(*newDb)

	log.Println("Starting Server...")
	server := router.Server{
		DB:           newDb,
		UsersUseCase: usersUseCase,
	}

	http.ListenAndServe(":3000", router.NewRouter(&server))
}
