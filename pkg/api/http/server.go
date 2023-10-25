package http

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"

	"github.com/tonytcb/party-invite/pkg/infrastructure/logger"
)

type Server struct {
	log        logger.Logger
	httpServer *http.Server

	filterCustomersHandler *FilterCustomersHandler
}

func NewServer(
	log logger.Logger,
	filterCustomersHandler *FilterCustomersHandler,
) *Server {
	return &Server{
		log:                    log,
		filterCustomersHandler: filterCustomersHandler,
	}
}

func (s *Server) Start(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.healthHandler)
	mux.HandleFunc("/filter-customers", s.filterCustomersHandler.Handle)

	s.log.Infof("Starting HTTP Server on port %d", port)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Errorf("error to listen and server http api: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Infof("Shutting down HTTP Server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "error to shutdown http server")
	}

	return nil
}

func (s *Server) healthHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "ok")
}
