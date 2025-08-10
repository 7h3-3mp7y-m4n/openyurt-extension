package handlers

import "sync"

type Status struct {
	Installed bool   `json:"installed"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var (
	currentStatus = Status{Installed: false, Status: "not_installed", Message: "OpenYurt not installed"}
	statusMutex   sync.RWMutex
)

func UpdateStatus(installed bool, status, message string) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	currentStatus = Status{Installed: installed, Status: status, Message: message}
}

func GetStatus() Status {
	statusMutex.RLock()
	defer statusMutex.RUnlock()
	return currentStatus
}
