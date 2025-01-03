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
	"github.com/Abo-Omar-74/httpServer/model"
	"github.com/google/uuid"
)

// CreateUserHandler handles the creation of a new user in the system.
func (h *Handler)CreateUserHandler(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodPost{
		helper.RespondWithError(w,http.StatusMethodNotAllowed , "Only POST requests are supported.")
		return
	}

	if h.Cfg.Platform != "dev"{
		helper.RespondWithError(w,http.StatusForbidden , "Access is allowed only in the development environment.")
		return
	}

	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var params parameters
	
	body , err := io.ReadAll(r.Body)
	if err != nil{
		helper.RespondWithError(w,http.StatusBadRequest , "Failed to read request body.")
		return 
	}

	err = json.Unmarshal(body , &params)

	if err != nil{
		helper.RespondWithError(w , http.StatusBadRequest ,"Invalid request format.")
		return
	}

	Hash , err := auth.HashPassword(params.Password)

	if err != nil{
		helper.RespondWithError(w , http.StatusInternalServerError ,"An error occurred while processing your request.")
		return
	}
	
	_ , err = h.Cfg.Db.FindUserByEmail(r.Context() ,params.Email)

	if err == nil{
		helper.RespondWithError(w,http.StatusConflict , "Email already exists.")
		return
	}else if err != sql.ErrNoRows{
		helper.RespondWithError(w,http.StatusInternalServerError , "An error occurred while processing your request.")
		return
	}

	dbUser , err := h.Cfg.Db.CreateUser(r.Context() , database.CreateUserParams{Email: params.Email , HashedPassword: Hash})
	if err != nil{
		helper.RespondWithError(w,http.StatusInternalServerError , "An error occurred while processing your request.")
		return
	}
	helper.RespondWithJSON(w,http.StatusCreated , model.DatabaseUserToUser(dbUser))
}


// EditUserHandler handles the request to edit a user's information.
func  (h *Handler) EditUserHandler(w http.ResponseWriter , r *http.Request , jwtUserID uuid.UUID){
	if  r.Method != http.MethodPut{
		helper.RespondWithError(w , http.StatusMethodNotAllowed , "Only PUT requests are allowed.")
		return	
	}
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var params parameters

	body , err := io.ReadAll(r.Body)
	if err != nil{
		helper.RespondWithError(w,http.StatusBadRequest , "Failed to read request body.")
		return 
	}
	err = json.Unmarshal(body , &params)
	if err != nil{
		helper.RespondWithError(w,http.StatusBadRequest , "Invalid request format.")
		return 
	}
	
	user , err := h.Cfg.Db.FindUserByID(r.Context() , jwtUserID)
	
	if err != nil {
		helper.RespondWithError(w,http.StatusUnauthorized ,  "Unauthorized: Invalid credentials.")
		return
	}

	if user.Email != params.Email{
		helper.RespondWithError(w , http.StatusUnauthorized ,  "Unauthorized: Invalid credentials.")
		return
	}
	
	HashedPassword , err := auth.HashPassword(params.Password)
	if err != nil {
		helper.RespondWithError(w,http.StatusInternalServerError , "An error occurred while processing the request.")
		return
	}
	newUser , err := h.Cfg.Db.EditUserByID(r.Context() , 
	database.EditUserByIDParams{
		ID: jwtUserID ,
		Email: params.Email ,
		HashedPassword: HashedPassword,
		})
	
	if err != nil {
		helper.RespondWithError(w,http.StatusInternalServerError , "An error occurred while processing the request.")
		return
	}
	newUser.UpdatedAt = time.Now()
	
	helper.RespondWithJSON(w , http.StatusOK , model.DatabaseUserToUser(newUser) )
}


// DeleteAllUsers handles the deletion of all users, restricted to the "dev" platform.
func (h *Handler) DeleteAllUsers(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodPost{
		helper.RespondWithError(w,http.StatusMethodNotAllowed , "Only POST requests are supported")
		return
	}
	if h.Cfg.Platform != "dev"{
		helper.RespondWithError(w,http.StatusForbidden , "Access is allowed only in the development environment.")
		return
	}
	
	err := h.Cfg.Db.DeleteAllUsers(r.Context())
	if err != nil{
		helper.RespondWithError(w,http.StatusInternalServerError , "Unable to process your request")
		return
	}
	helper.RespondWithJSON(w,http.StatusOK,"All Users have been deleted successfully")
}