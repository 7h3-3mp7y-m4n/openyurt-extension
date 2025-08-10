package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/7h3-3mp7y-m4n/open-extension/backend/utils"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w)
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	status := GetStatus()
	response := Response{Success: true, Message: "Status retrieved", Data: status}
	json.NewEncoder(w).Encode(response)
}
