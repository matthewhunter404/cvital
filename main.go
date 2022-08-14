package main

import (
	"cvital/db"
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
	err = db.RunMigrations(newDb)
	if err != nil {
		log.Fatalf("DB migrations failed: %v", err)
	}
	log.Println("Starting Server...")
	http.ListenAndServe(":3000", router.NewRouter())
}
