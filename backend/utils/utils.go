package utils

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func EnableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

var root string

func RunScript(scriptName string, args ...string) (string, error) {
	if root == "" {
		wd, _ := os.Getwd()
		for _, r := range []string{wd, filepath.Join(wd, ".."), filepath.Join(wd, "..", "..")} {
			if _, err := os.Stat(filepath.Join(r, "scripts")); err == nil {
				root, _ = filepath.Abs(r)
				break
			}
		}
		if root == "" {
			return "", fmt.Errorf("scripts directory not found")
		}
	}

	script := filepath.Join(root, "scripts", scriptName)
	os.Chmod(script, 0755)
	cmd := exec.Command(script, args...)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	return string(output), err
}
