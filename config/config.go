package config

import "github.com/Abo-Omar-74/httpServer/internal/database"


type ApiConfig struct{
  Db *database.Queries
  Platform string
  JwtSecret string
  UpgradePremiumKey string
}