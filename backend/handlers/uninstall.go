package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/7h3-3mp7y-m4n/open-extension/backend/utils"
)

func UninstallHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w)
	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	currentStatus := GetStatus()
	if !currentStatus.Installed {
		response := Response{Success: false, Message: "OpenYurt is not installed"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if currentStatus.Status == "uninstalling" {
		response := Response{Success: false, Message: "Uninstallation already in progress"}
		json.NewEncoder(w).Encode(response)
		return
	}

	go func() {
		UpdateStatus(true, "uninstalling", "Removing OpenYurt components...")

		output, err := utils.RunScript("uninstall-script.sh")
		if err != nil {
			log.Printf("Uninstallation failed: %v, Output: %s", err, output)
			UpdateStatus(true, "failed", "Uninstallation failed: "+err.Error())
			return
		}

		UpdateStatus(false, "not_installed", "OpenYurt removed successfully")
		log.Println("OpenYurt uninstallation completed")
	}()

	response := Response{Success: true, Message: "Uninstallation started"}
	json.NewEncoder(w).Encode(response)
}
