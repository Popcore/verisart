package handlers

import (
	"fmt"
	"net/http"
)

// GetHome returns
func GetHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello server")
}
