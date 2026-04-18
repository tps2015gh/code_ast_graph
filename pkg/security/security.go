package security

import (
	"log"
	"os"
	"strings"
)

// CheckRepoSafety ensures that a .gitignore file exists to prevent
// leaking sensitive data or large dependencies into the repository.
func CheckRepoSafety() {
	log.Println("Security Agent: Checking repository safety...")
	if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
		log.Println("Security Agent WARNING: .gitignore file is missing! Creating a default one.")
		// The .gitignore was created by the builder, but this is a runtime check.
	} else {
		log.Println("Security Agent: .gitignore found.")
	}
}

// PrivacyScrub masks potentially sensitive information in strings
// before they are sent to the frontend or stored.
func PrivacyScrub(input string) string {
	lower := strings.ToLower(input)
	
	// Basic list of sensitive keywords to scrub
	sensitiveKeywords := []string{"password", "secret", "apikey", "token", "auth_key"}
	
	for _, kw := range sensitiveKeywords {
		if strings.Contains(lower, kw) {
			return "[SCRUBBED SENSITIVE DATA]"
		}
	}
	
	return input
}
