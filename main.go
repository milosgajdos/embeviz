package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	v1 "github.com/milosgajdos/embeviz/api/v1"
	"github.com/milosgajdos/embeviz/api/v1/http"
	"github.com/milosgajdos/embeviz/api/v1/memory"
	"github.com/milosgajdos/go-embeddings/cohere"
	"github.com/milosgajdos/go-embeddings/openai"
	"github.com/milosgajdos/go-embeddings/vertexai"
	"golang.org/x/oauth2/google"
)

const (
	// cliName is command line name.
	cliName = "embeviz"

	// http.Server timeouts
	IdleTimeout  = 5 * time.Second
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 10 * time.Second

	// ShutdownTimeout defines time when we forcefully shutdown the server
	ShutdownTimeout = 10 * time.Second
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(cliName, flag.ExitOnError)

	var (
		addr = flags.String("addr", ":5050", "API server bind address")
		dsn  = flags.String("dsn", ":memory:", "Database connection string")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	options := []http.Option{
		http.WithIdleTimeout(IdleTimeout),
		http.WithReadTimeout(ReadTimeout),
		http.WithWriteTimeout(WriteTimeout),
	}

	s, err := http.NewServer(options...)
	if err != nil {
		return err
	}

	// creates provider service
	ps, err := NewProviderService(*dsn)
	if err != nil {
		return err
	}

	// adds default embedders
	embedders, err := AddDefaultEmbedders(ps)
	if err != nil {
		return err
	}

	s.Addr = *addr
	s.ProvidersService = ps
	s.Embedders = embedders

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error, 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- s.Listen()
	}()

	// Listen for the interrupt signal.
	select {
	case <-ctx.Done():
	case err := <-errChan:
		return err
	}

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// Perform application shutdown with a maximum timeout of ShutdownTimeout seconds.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	return s.Close(timeoutCtx)
}

// NewProviderService returns an instance of v1.ProviderService based on the DSN
// Currently we only support in-memory implementation
func NewProviderService(dsn string) (v1.ProvidersService, error) {
	switch dsn {
	case memory.DSN:
		return makeMemoryProvidersService(dsn)
	default:
		return nil, fmt.Errorf("unsuported DSN: %s", dsn)
	}
}

// makeMemoryProvidersService creates an ProvidersService
// backed my in-memory datastore and returns it.
func makeMemoryProvidersService(dsn string) (v1.ProvidersService, error) {
	db, err := memory.NewDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed creating new DB: %v", err)
	}
	if err := db.Open(); err != nil {
		return nil, fmt.Errorf("failed opening DB: %v", err)
	}
	return memory.NewProvidersService(db)
}

// AddDefaultEmbedders adds the default embedders as long as alll the required
// environment variables are set for eache specific embedding provider:
// * OPENAI_API_KEY for OpenAI API
// * COHERE_API_KEY for Cohere API
// * VERTEXAI_TOKEN, VERTEXAI_MODEL_ID, GOOGLE_PROJECT_ID for Google VertexAI
func AddDefaultEmbedders(ps v1.ProvidersService) (map[string]any, error) {
	embedders := make(map[string]any)

	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		openAI, err := ps.AddProvider(context.Background(), "OpenAI", map[string]any{})
		if err != nil {
			return nil, err
		}
		embedders[openAI.UID] = openai.NewClient()
	}

	if apiKey := os.Getenv("COHERE_API_KEY"); apiKey != "" {
		cohereAI, err := ps.AddProvider(context.Background(), "Cohere", map[string]any{})
		if err != nil {
			return nil, err
		}
		embedders[cohereAI.UID] = cohere.NewClient()
	}

	if os.Getenv("VERTEXAI_TOKEN") != "" &&
		os.Getenv("GOOGLE_PROJECT_ID") != "" {
		vertexAI, err := ps.AddProvider(context.Background(), "VertexAI", map[string]any{})
		if err != nil {
			return nil, err
		}

		ts, err := google.DefaultTokenSource(context.Background(), vertexai.Scopes)
		if err != nil {
			return nil, fmt.Errorf("vertexai: token source: %v", err)
		}
		embedders[vertexAI.UID] = vertexai.NewClient(
			vertexai.WithTokenSrc(ts),
			vertexai.WithModelID(vertexai.EmbedGeckoV2.String()))
	}

	return embedders, nil
}
