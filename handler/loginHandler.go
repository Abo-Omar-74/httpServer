package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Abo-Omar-74/httpServer/helper"
	"github.com/Abo-Omar-74/httpServer/internal/auth"
	"github.com/Abo-Omar-74/httpServer/internal/database"
	"github.com/google/uuid"
)




type LoginResponse struct{
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
  Email string `json:"email"` 
	AccessToken string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsPremium bool `json:"is_premium"`
}

type RefreshResponse struct{
	RefreshToken string `json:"token"`
}

// LoginHandler handles user login by verifying the email and password.
// If authentication is successful, it returns an authentication token and a refresh token.
func (h *Handler) LoginHandler(w http.ResponseWriter , r *http.Request){
	
	type parameters struct{
		Password string `json:"password"`
		Email string `json:"email"`
	}

	if r.Method != http.MethodPost{
		helper.RespondWithError(w,http.StatusMethodNotAllowed , "Method Not Allowed !")
		return
	}
	var params parameters
	body , err := io.ReadAll(r.Body)
	if err != nil{
		helper.RespondWithError(w,http.StatusBadRequest , "Invalid request body.")
		return 
	}
	err = json.Unmarshal(body , &params)
	if err != nil{
		helper.RespondWithError(w,http.StatusBadRequest , "Invalid JSON format.")
		return 
	}

	user , err := h.Cfg.Db.FindUserByEmail(r.Context() , params.Email)

	if err != nil{
		helper.RespondWithError(w,http.StatusUnauthorized , "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword , params.Password)
	if err != nil{
		helper.RespondWithError(w,http.StatusUnauthorized , "Incorrect email or password")
		return
	}


	token , err:= auth.MakeJWT(user.ID , h.Cfg.JwtSecret)

	if err != nil{
		helper.RespondWithError(w,http.StatusInternalServerError , "An error occurred while processing your request.")
		return
	}
	refreshToken , err := auth.MakeRefreshToken()

	if err != nil{
		helper.RespondWithError(w , http.StatusInternalServerError , "An error occurred while processing your request.")
		return
	}
	res := LoginResponse{user.ID , user.CreatedAt , user.UpdatedAt , user.Email , token , refreshToken , user.IsPremium}
	helper.RespondWithJSON(w,http.StatusAccepted,res)
}

// RefreshHandler handles the refreshing of JWT tokens using a valid refresh token.
func (h *Handler) RefreshHandler(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodPost{
		helper.RespondWithError(w,http.StatusMethodNotAllowed , "Only POST requests are supported.")
		return
	}

	reqRefreshToken , err := auth.GetBearerToken(r.Header)
	if err != nil{
		helper.RespondWithError(w , http.StatusUnauthorized , "Invalid or missing refresh token.")
		return
	}

	dbRefreshToken , err := h.Cfg.Db.GetRefreshToken(r.Context() , reqRefreshToken)
	if err != nil{
		helper.RespondWithError(w , http.StatusUnauthorized , "Invalid refresh token.")
		return
	}
	if dbRefreshToken.ExpiresAt.Before(time.Now()) || dbRefreshToken.RevokedAt.Valid{
		helper.RespondWithError(w , http.StatusUnauthorized , "Refresh token is no longer valid.")
		return
	} 
	jwtToken , err := auth.MakeJWT(dbRefreshToken.UserID , h.Cfg.JwtSecret)
	
	if err != nil{
		helper.RespondWithError(w , http.StatusInternalServerError , "Failed to generate token.")
		return
	}

	helper.RespondWithJSON(w,http.StatusAccepted , RefreshResponse{jwtToken})
}

// RevokeHandler revokes a user's refresh token, making it invalid for future use.
func (h *Handler) RevokeHandler(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodPost{
		helper.RespondWithError(w,http.StatusMethodNotAllowed , "Only POST requests are supported.")
		return
	}

	reqRefreshToken , err := auth.GetBearerToken(r.Header)
	if err != nil{
		helper.RespondWithError(w , http.StatusUnauthorized , "Invalid or missing refresh token.")
		return
	}

	refreshToken , err := h.Cfg.Db.GetRefreshToken(r.Context() , reqRefreshToken)
	if err != nil{
		helper.RespondWithError(w , http.StatusUnauthorized , "Invalid refresh token.")
		return
	}

	refreshToken.RevokedAt = sql.NullTime{
		Time : time.Now(), 
		Valid: true,
	}

	err = h.Cfg.Db.RevokeRefreshToken(r.Context() , database.RevokeRefreshTokenParams{
			RevokedAt : refreshToken.RevokedAt,
			UpdatedAt : time.Now(),
			Token     : refreshToken.Token,
	})

	if err != nil{
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token.")
		return
	}

	helper.RespondWithJSON(w,http.StatusNoContent , nil)
}