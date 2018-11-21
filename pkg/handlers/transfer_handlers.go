package handlers

import (
	"fmt"
	"net/http"
)

// GetHome returns
func PostTransferHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello server")
}

func PatchTransferHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello server")
}
