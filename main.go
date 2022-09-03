package main

import (
	"cvital/config"
	"cvital/db"
	"cvital/domain/profiles"
	"cvital/domain/users"
	"cvital/router"
	"fmt"
	"log"
	"net/http"
)

func main() {

	config, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Reading config file failed: %v", err)
	}

	newDb, err := db.NewConnection(db.DatabaseConfig{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		User:     config.Database.User,
		DbName:   config.Database.Name,
		Password: config.Database.Password,
		SslMode:  config.Database.SslMode,
	})
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	err = db.RunMigrations(newDb.DB)
	if err != nil {
		log.Fatalf("DB migrations failed: %v", err)
	}

	usersUseCase := users.NewUseCase(*newDb, config.JWTKey)
	profilesUseCase := profiles.NewUseCase(*newDb)
	log.Println("Starting Server...")
	server := router.Server{
		DB:              newDb,
		UsersUseCase:    usersUseCase,
		ProfilesUseCase: profilesUseCase,
	}

	http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router.NewRouter(&server))
}
