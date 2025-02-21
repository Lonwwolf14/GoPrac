package app

import (
	"httpServer/internal/config"
	"httpServer/internal/database"
)

type AppState struct {
	AppConfig *config.Config
	DB        *database.Queries
}
