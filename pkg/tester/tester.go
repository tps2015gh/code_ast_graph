package tester

import (
	"log"
	"net/http"
	"net/http/httptest"
)

// RunBasicTests performs a simple check on the provided HTTP handler.
// This is the starting point for the "tester agent".
func RunBasicTests(handler http.Handler) {
	log.Println("Tester Agent: Starting basic health checks...")

	// Test a simple request to the root
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Printf("Tester Agent: Failed to create request: %v", err)
		return
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		log.Printf("Tester Agent: Health check failed: expected status 200, got %d", status)
	} else {
		log.Println("Tester Agent: Health check passed (root endpoint is accessible).")
	}

	log.Println("Tester Agent: Basic health checks complete.")
}
