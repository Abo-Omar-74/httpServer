package helper

import (
	"encoding/json"
	"net/http"
)
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error{
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	
	// Avoid writing body for 204 no content status
	if code == http.StatusNoContent{
		return nil
	}
	
	response , err := json.Marshal(payload)
	if(err != nil){
		return err
	}

	w.Write(response)
	return nil
}
func RespondWithError(w http.ResponseWriter, code int, msg string) error{

	return RespondWithJSON(w , code , map[string]string{"error":msg})

}
