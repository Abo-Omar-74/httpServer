package middleware

import (
	"log"
	"net/http"

	"github.com/Abo-Omar-74/httpServer/config"
	"github.com/Abo-Omar-74/httpServer/helper"
	"github.com/Abo-Omar-74/httpServer/internal/auth"
	"github.com/google/uuid"
)

type Middleware struct {
	Cfg *config.ApiConfig
}


type authedHandler func(w http.ResponseWriter , r *http.Request , jwtuserID uuid.UUID )

func (m *Middleware) MiddlewareAuth(handler authedHandler) http.HandlerFunc{
	return func (w http.ResponseWriter , r *http.Request){

		token , err := auth.GetBearerToken(r.Header)

		if err != nil{
			log.Printf("Missing or invalid token: %v", err)
			helper.RespondWithError(w,http.StatusUnauthorized , "Unauthorized")
			return
		}
		userID , err := auth.ValidateJWT(token , m.Cfg.JwtSecret)

		if err != nil{
			log.Printf("Invalid JWT: %v", err)
			helper.RespondWithError(w,http.StatusUnauthorized ,"Unauthorized")
			return
		}
		handler(w , r , userID)
	}
}