package handlers

import (
	"fmt"
	"net/http"
)

func ListUserCertsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello server")
}
