package http

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	v1 "github.com/milosgajdos/embeviz/api/v1"
	_ "github.com/milosgajdos/embeviz/api/v1/http/docs" // blank import for swagger docs
)

// Server is an HTTP server used to provide REST API
// access for various Graph API endpoints.
type Server struct {
	// app is fiber app.
	app *fiber.App
	// ln is a network listener.
	ln net.Listener
	// Addr is bind address
	Addr string
	// ProvidersService provides access to Provider enpoints.
	ProvidersService v1.ProvidersService
	// Embedders
	// NOTE: this is a major hack
	Embedders map[string]any
}

type Options struct {
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Option is functional graph option.
type Option func(*Options)

// WithIdleTimeout
func WithIdleTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.IdleTimeout = t
	}
}

// WithReadTimeout
func WithReadTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = t
	}
}

// WithWriteTimeout
func WithWriteTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.WriteTimeout = t
	}
}

// NewServer creates a new Server and returns it.
func NewServer(options ...Option) (*Server, error) {
	var c fiber.Config

	opts := Options{}
	for _, apply := range options {
		apply(&opts)
	}

	c.IdleTimeout = opts.IdleTimeout
	c.ReadTimeout = opts.ReadTimeout
	c.WriteTimeout = opts.WriteTimeout

	s := &Server{
		app: fiber.New(c),
	}

	s.app.Use(recover.New())
	s.app.Use(logger.New())
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	api := s.app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/docs/*", swagger.New())

	s.registerProviderRoutes(v1)

	return s, nil
}

// Listen validates the server options and binds to the given address.
func (s *Server) Listen() error {
	if s.Addr == "" {
		return fmt.Errorf("empty bind address")
	}

	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.ln = ln

	return s.app.Listener(ln)
}

// Close gracefully shuts down the server.
func (s *Server) Close(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func() {
		select {
		case <-ctx.Done():
		case errChan <- s.app.Shutdown():
		}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("server shut down: %v", ctx.Err())
	case err := <-errChan:
		return err
	}
}
