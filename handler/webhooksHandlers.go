package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Abo-Omar-74/httpServer/helper"
	"github.com/Abo-Omar-74/httpServer/internal/auth"
	"github.com/google/uuid"
)

const UserUpgradedEvent = "user.upgraded"


// UpgradeUserHandler handles upgrading a user to premium status based on event data.
func (h *Handler)UpgradeUserHandler(w http.ResponseWriter , r *http.Request){
	
	type eventData struct { 
		UserID uuid.UUID `json:"user_id"` 
	} 
	
	type parameters struct { 
		Event string `json:"event"` 
		Data eventData `json:"data"` 
	}
	if r.Method != http.MethodPost{
		helper.RespondWithError(w , http.StatusMethodNotAllowed , "Only POST requests are supported")
		return
	}
	apiKey , err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != h.Cfg.UpgradePremiumKey{
		helper.RespondWithError(w,http.StatusUnauthorized , "Access denied")
		return
	}
	body , err := io.ReadAll(r.Body)
	if err != nil{
		helper.RespondWithError(w , http.StatusBadRequest , "Invalid Json Format")
		return
	}
	var params parameters
	err = json.Unmarshal(body , &params)
	if err != nil{
		helper.RespondWithError(w , http.StatusBadRequest , "Invalid JSON format")
		return
	}
	if params.Event != UserUpgradedEvent{
		helper.RespondWithJSON(w , http.StatusNoContent , nil)
		return
	}
	_ , err = h.Cfg.Db.UpgradeUserByID(r.Context() , params.Data.UserID)
	if err != nil{
		if errors.Is(err , sql.ErrNoRows){
				helper.RespondWithError(w , http.StatusNotFound , "User not found")
		}else {
				helper.RespondWithError(w , http.StatusInternalServerError , "Unable to process the request")
		}
		return
	}
	helper.RespondWithJSON(w , http.StatusNoContent , nil)
}