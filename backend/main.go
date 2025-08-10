package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/7h3-3mp7y-m4n/open-extension/backend/handlers"
)

func main() {
	http.HandleFunc("/status", handlers.StatusHandler)
	http.HandleFunc("/install", handlers.InstallHandler)
	http.HandleFunc("/uninstall", handlers.UninstallHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	fs := http.FileServer(http.Dir(filepath.Join("ui")))
	http.Handle("/", fs)
	fmt.Println("Thanks for looking at my code 😊")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Server can't be started --- %v \n", err)
	}
}
