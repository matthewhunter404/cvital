package main

import (
	"cvital/db"
	"cvital/router"
	"log"
	"net/http"
)

func main() {
	_, err := db.NewConnection(db.DatabaseConfig{
		Host:     "cvital",
		Port:     "5432",
		User:     "cvital",
		DbName:   "cvital",
		Password: "cvital",
		SslMode:  "false",
	})
	if err != nil {
		log.Fatalf("DB connection failed")
	}
	log.Println("Starting Server...")
	http.ListenAndServe(":3000", router.NewRouter())
}
