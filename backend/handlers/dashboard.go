package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/7h3-3mp7y-m4n/open-extension/backend/utils"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w)
	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	status := GetStatus()
	if !status.Installed {
		response := Response{Success: false, Message: "OpenYurt not installed"}
		json.NewEncoder(w).Encode(response)
		return
	}

	output, err := utils.RunScript("get-dashboard-url.sh") // I didn't add this yet, TODO I guess ;)
	if err != nil {
		response := Response{Success: false, Message: "Dashboard not available"}
		json.NewEncoder(w).Encode(response)
		return
	}

	dashboardURL := strings.TrimSpace(output)
	response := Response{Success: true, Message: "Dashboard available", Data: map[string]string{"url": dashboardURL}}
	json.NewEncoder(w).Encode(response)
}
