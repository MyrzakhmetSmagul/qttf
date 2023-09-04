package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"qttf/config"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	cnf  *config.Config
	db   *sql.DB
}

func NewServer(cnf *config.Config, db *sql.DB) *Server {
	return &Server{echo: echo.New(), cnf: cnf, db: db}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:         s.cnf.Router.Port,
		ReadTimeout:  time.Duration(s.cnf.Router.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.cnf.Router.WriteTimeout) * time.Second,
	}

	go func() {
		log.Printf("Server is listening on port: %s\n", s.cnf.Router.Port)
		if err := s.echo.StartServer(server); err != nil {
			log.Fatalf("can't start server %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	log.Println("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
