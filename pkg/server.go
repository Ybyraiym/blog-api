package blog

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,       //адрес, на котором сервер будет слушать входящие соединения.
		MaxHeaderBytes: 1 << 20,          //1 MB
		ReadTimeout:    10 * time.Second, // Тайм-ауты на чтение и запись данных для каждого соединения.
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
