package server

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"avoxi-api/routes"
)

// Server is the core structure for handling web requests
type Server struct {
	port string
}

type errorResponse struct {
	HTTPStatusCode int `json:"httpStatusCode"`
	ErrorMessage string `json:"errorMessage"`
}

// Start sets up the routes for the HTTP server and defines the port to listen on then handles all inbound requests
func (s *Server) Start(port int) error {
	s.setPort(port)
	s.setupRoutes()

	err := s.listenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) setPort(port int) {
	s.port = fmt.Sprintf(":%v", port)
}

func (s *Server) setupRoutes() {
	routes.Setup(s.handleError)
}

func (s *Server) listenAndServe() error {
	err := http.ListenAndServe(s.port, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleError(status int, err error, w http.ResponseWriter, r *http.Request) {
	response := &errorResponse{
		HTTPStatusCode: status,
	}

	w.WriteHeader(status)

	switch status {
	case http.StatusBadRequest:
		log.Printf("Failed to parse inbound request. Error Details: %v", err)
		response.ErrorMessage = "Failed to parse request. Please ensure request JSON is valid and try your request again."
	case http.StatusInternalServerError:
		log.Printf("An unexpected error occurred. Error details: %v", err)
		response.ErrorMessage = "An unexpected error occurred. Please wait some time and try your request again."
	case http.StatusNotImplemented:
		log.Printf("The requested endpoint %s does not have a way to process %s requests.", r.RequestURI, r.Method)
		response.ErrorMessage = "That request type is not available for this endpoint. Please select an available request type and try again."
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal error response after an error was encountered. Details: %v", err)
	} else {
		w.Write(bytes)
	}
}