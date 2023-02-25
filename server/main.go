package main

import (
	"cvital/config"
	"cvital/db"
	"cvital/domain/profiles"
	"cvital/domain/users"
	"cvital/logger"
	"cvital/router"
	"fmt"
	"net/http"
)

func main() {
	logger := logger.NewLogger()
	config, err := config.ReadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Reading config file failed")
	}

	newDb, err := db.NewConnection(db.DatabaseConfig{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		User:     config.Database.User,
		DbName:   config.Database.Name,
		Password: config.Database.Password,
		SslMode:  config.Database.SslMode,
	}, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("DB connection failed")
	}

	usersUseCase := users.NewUseCase(newDb, config.JWTKey, logger)
	profilesUseCase := profiles.NewUseCase(newDb, logger)
	logger.Info().Msg("Starting Server...")
	server := router.Server{
		DB:              newDb,
		UsersUseCase:    usersUseCase,
		ProfilesUseCase: profilesUseCase,
		Logger:          logger,
	}

	http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router.NewRouter(&server))
}
