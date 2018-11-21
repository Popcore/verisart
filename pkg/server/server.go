package server

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"goji.io"
	"goji.io/pat"

	"github.com/popcore/verisart_exercise/pkg/certificate"
	cert "github.com/popcore/verisart_exercise/pkg/certificate"
	"github.com/popcore/verisart_exercise/pkg/handlers"
)

// Server is a custom type used to group server configuration,
// services and functionalities
type Server struct {
	Address   string
	Mux       *goji.Mux
	CertStore certificate.Storer
}

// New returns a server instance than can be used to handle
// http requests.
func New(addr string) *Server {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/certificates:id"), handlers.GetHome)
	mux.HandleFunc(pat.Patch("/certificates:id"), handlers.GetHome)
	mux.HandleFunc(pat.Delete("/certificates:id"), handlers.GetHome)
	mux.HandleFunc(pat.Get("/users/:userId/certificates"), handlers.GetHome)
	mux.HandleFunc(pat.Post("/users/:userId/transfers"), handlers.GetHome)
	mux.HandleFunc(pat.Patch("/users/:userId/transfers"), handlers.GetHome)

	// set up cors
	c := cors.New(
		cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT"},
			AllowedHeaders: []string{"Authorization", "Content-Type"},
		},
	)
	mux.Use(c.Handler)

	return &Server{
		Address:   addr,
		Mux:       mux,
		CertStore: cert.NewMemStore(),
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
