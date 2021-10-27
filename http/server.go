package http

import (
	"net/http"
)

//
type Server struct {
	server *http.Server
}

//
func NewServer() *Server {
	srv := &Server{
		server: &http.Server{
			Addr:              "",
			Handler:           nil,
			TLSConfig:         nil,
			ReadTimeout:       0,
			ReadHeaderTimeout: 0,
			WriteTimeout:      0,
			IdleTimeout:       0,
			MaxHeaderBytes:    0,
			TLSNextProto:      nil,
			ConnState:         nil,
			ErrorLog:          nil,
			BaseContext:       nil,
			ConnContext:       nil,
		},
	}

	return srv
}

func (srv *Server) Start() error {
	return nil
}

func (srv *Server) Stop() {

}
