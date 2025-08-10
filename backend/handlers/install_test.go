package handlers

import (
	"log"
	"testing"

	"github.com/7h3-3mp7y-m4n/open-extension/backend/utils"
)

func TestRunScript(t *testing.T) {
	output, err := utils.RunScript("01-check-prerequisites.sh")
	if err != nil {
		log.Printf("Passed : %v, Output: %s", err, output)
		UpdateStatus(false, "failed", "Pre-Req failed: "+err.Error())
		t.Fatalf("Script failed with error: %v", err)

	}

	t.Logf("Script ran successfully. Output:\n%s", output)
}
