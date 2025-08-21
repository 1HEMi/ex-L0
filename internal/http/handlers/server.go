package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/orders").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateOrder)
	router.Path("/orders/{id}").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetOrder)
	router.Path("/orders").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetOrders)
	router.Path("/orders/{id}").Methods("DELETE").HandlerFunc(s.httpHandlers.HandleDeleteOrder)

	return http.ListenAndServe(":8080", router)
}
