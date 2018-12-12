package server

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"goji.io"
	"goji.io/pat"

	"github.com/Popcore/verisart/pkg/handlers"
	store "github.com/Popcore/verisart/pkg/store"
)

// Server is a custom type used to group server configuration,
// services and functionalities
type Server struct {
	Address string
	Mux     *goji.Mux
}

// New returns a server instance than can be used to handle
// http requests.
func New(addr string) *Server {

	memStore := store.NewMemStore()
	mux := goji.NewMux()
	mux.Handle(pat.Post("/certificates"), handlers.Handler{S: memStore, H: handlers.PostCertHandler})
	mux.Handle(pat.Patch("/certificates/:id"), handlers.Handler{S: memStore, H: handlers.PatchCertHandler})
	mux.Handle(pat.Delete("/certificates/:id"), handlers.Handler{S: memStore, H: handlers.DeleteCertHandler})
	mux.Handle(pat.Post("/certificates/:id/transfers"), handlers.Handler{S: memStore, H: handlers.PostTransferHandler})
	mux.Handle(pat.Patch("/certificates/:id/transfers"), handlers.Handler{S: memStore, H: handlers.PatchTransferHandler})
	mux.Handle(pat.Get("/users/:userId/certificates"), handlers.Handler{S: memStore, H: handlers.ListUserCertsHandler})
	mux.Handle(pat.Post("/users"), handlers.Handler{S: memStore, H: handlers.NewUserHandler})
	// define cors policies
	c := cors.New(
		cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
			AllowedHeaders: []string{"Authorization", "Content-Type", "X-User-Email"},
		},
	)
	mux.Use(c.Handler)

	return &Server{
		Address: addr,
		Mux:     mux,
	}
}

// Start generates and runs and http.Server at the defined address.
func (s *Server) Start() {
	servMux := &http.Server{
		Addr:    s.Address,
		Handler: s.Mux,
	}

	log.Printf("Server running at %s", s.Address)

	err := servMux.ListenAndServe()
	if err != nil {
		log.Fatalf("Unexpected error starting http server: %s", err.Error())
	}
}
