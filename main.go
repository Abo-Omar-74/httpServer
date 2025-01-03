package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Abo-Omar-74/httpServer/config"
	"github.com/Abo-Omar-74/httpServer/handler"
	"github.com/Abo-Omar-74/httpServer/internal/database"
	"github.com/Abo-Omar-74/httpServer/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
  port = "8080"
  filepathRoot = "."
)



func main(){
  godotenv.Load()

  jwtSecret := os.Getenv("JWT_SECRET")
  if jwtSecret == ""{
    log.Fatal("JWT_SECRET is not set")
  }
  upgradePremiumKey := os.Getenv("API_KEY_UPGRADE_PREMIUM")
  if upgradePremiumKey == ""{
    log.Fatal("API_KEY_UPGRADE_PREMIUM is not set")
  }

  dbURL := os.Getenv("DB_URL")
  db, err := sql.Open("postgres", dbURL)
  
  if err != nil{
    fmt.Print(err)
    return
  }
  dbQueries := database.New(db)
  platform := os.Getenv("PLATFORM")

  apiCfg := config.ApiConfig{
    Db : dbQueries,
    Platform: platform,
    JwtSecret: jwtSecret,
    UpgradePremiumKey: upgradePremiumKey,
  }

  apiHandler := &handler.Handler{
    Cfg: &apiCfg,
  }
  apiMiddleware := &middleware.Middleware{
    Cfg : &apiCfg,
  }

  mux := http.NewServeMux()

  // Pattern - Handlers Binding

  mux.HandleFunc("POST /api/users" ,  apiHandler.CreateUserHandler)
  mux.HandleFunc("PUT /api/users" ,   apiMiddleware.MiddlewareAuth(apiHandler.EditUserHandler))
  mux.HandleFunc("POST /admin/reset", apiHandler.DeleteAllUsers)
  
  mux.HandleFunc("POST /api/login" ,  apiHandler.LoginHandler)
  mux.HandleFunc("POST /api/refresh", apiHandler.RefreshHandler)
  mux.HandleFunc("POST /api/revoke",  apiHandler.RevokeHandler)

  
  mux.HandleFunc("POST /api/posts" , apiMiddleware.MiddlewareAuth(apiHandler.PostHandler))
  mux.HandleFunc("GET /api/posts"  , apiHandler.GetPostsHandler)
  mux.HandleFunc("GET /api/posts/{postID}"    , apiHandler.GetPostByIDHandler)
  mux.HandleFunc("DELETE /api/posts/{postID}" , apiMiddleware.MiddlewareAuth(apiHandler.DeletePostHandler))


  mux.HandleFunc("POST /api/upgrade-premium/webhooks" , apiHandler.UpgradeUserHandler)



  // Create http serve to handel incoming request with patterns set before

  server := &http.Server{Addr : ":" + port , Handler : mux}

  log.Printf("Serving files %s on port: %s\n" , filepathRoot , port)

  log.Fatal(server.ListenAndServe())

}