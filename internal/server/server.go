package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/helper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Host         string
	Port         string
	Handler      http.Handler
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ErrorLog     *log.Logger
}

type Options func(*Server)

func WithHost(host string) Options {
	return func(s *Server) {
		s.Host = host
	}
}

func WithPort(port string) Options {
	return func(s *Server) {
		s.Port = port
	}
}

func WithHandler(handler http.Handler) Options {
	return func(s *Server) {
		s.Handler = handler
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Options {
	return func(s *Server) {
		s.IdleTimeout = idleTimeout
	}
}

func WithReadTimeout(readTimeout time.Duration) Options {
	return func(s *Server) {
		s.ReadTimeout = readTimeout
	}
}

func WithWriteTimeout(writeTimeout time.Duration) Options {
	return func(s *Server) {
		s.WriteTimeout = writeTimeout
	}
}

func WithErrorLog(errorLog *log.Logger) Options {
	return func(s *Server) {
		s.ErrorLog = errorLog
	}
}

func (s *Server) Connect() error {
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      s.Handler,
		IdleTimeout:  s.IdleTimeout,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		ErrorLog:     s.ErrorLog,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		se := <-quit
		log.Println("caught signal", "signal", se.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		log.Println("completing background tasks", "addr", server.Addr)
		helper.WG.Wait()
		shutdownError <- nil
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdownError; err != nil {
		return err
	}

	log.Println("Stopped server", "addr", server.Addr)

	return nil
}

func NewServer(opts ...Options) *Server {
	server := &Server{}
	for _, opt := range opts {
		opt(server)
	}
	return server
}
