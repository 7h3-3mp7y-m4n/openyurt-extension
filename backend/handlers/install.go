package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/7h3-3mp7y-m4n/open-extension/backend/utils"
)

func InstallHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w)
	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	go func() {
		UpdateStatus(false, "installing", "Checking prerequisites...")
		output, err := utils.RunScript("01-check-prerequisites.sh")
		if err != nil {
			log.Printf("Prerequisites check failed: %v, Output: %s", err, output)
			UpdateStatus(false, "failed", "Prerequisites check failed: "+err.Error())
			return
		}
		log.Printf("Prerequisites check passed: %s", output)
		UpdateStatus(false, "installing", "Setting up OpenYurt Helm repository...")
		output, err = utils.RunScript("02-setup-helm-repo.sh")
		if err != nil {
			log.Printf("Helm repo setup failed: %v, Output: %s", err, output)
			UpdateStatus(false, "failed", "Helm repo setup failed: "+err.Error())
			return
		}
		log.Printf("Helm repo setup completed: %s", output)

		UpdateStatus(false, "installing", "Installing OpenYurt control plane components...")
		output, err = utils.RunScript("03-install-yurt-manager.sh")
		if err != nil {
			log.Printf("OpenYurt installation failed: %v, Output: %s", err, output)
			UpdateStatus(false, "failed", "OpenYurt installation failed: "+err.Error())
			return
		}
		log.Printf("OpenYurt installation completed: %s", output)

		UpdateStatus(false, "installing", "Verifying OpenYurt installation...")
		output, err = utils.RunScript("04-verify-installation.sh")
		if err != nil {
			log.Printf("Installation verification failed: %v, Output: %s", err, output)
			UpdateStatus(false, "failed", "Installation verification failed: "+err.Error())
			return
		}
		log.Printf("Installation verification passed: %s", output)

		UpdateStatus(true, "running", "OpenYurt installed successfully")
		log.Println("OpenYurt installation completed successfully")
	}()

	response := Response{Success: true, Message: "Installation started"}
	json.NewEncoder(w).Encode(response)
}
