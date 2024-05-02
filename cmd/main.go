package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"bitswan.space/container-discovery-service/internal/config"
	"bitswan.space/container-discovery-service/internal/http"
	"bitswan.space/container-discovery-service/internal/logger"
	"bitswan.space/container-discovery-service/internal/mqtt"
	"bitswan.space/container-discovery-service/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	// Define a command-line flag
	configPath := flag.String("c", "config.yaml", "path to the configuration file")
	flag.Parse() // Parse the flags

	logger.Init()
	godotenv.Load(".env")

	err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error.Fatalf("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	err = mqtt.Init()
	if err != nil {
		logger.Error.Fatalf("Failed to initialize MQTT client: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Instantiate a new type to represent our application.
	// This type lets us shared setup code with our end-to-end tests.
	m := NewApp()

	if err := m.RunHTTPServer(ctx); err != nil {
		// m.Close()
		logger.Error.Printf("failed to start because %v", err)
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-sigChan
	logger.Info.Println("Shutting down gracefully...")
	// Perform any necessary cleanup here
	mqtt.Close()

	logger.Info.Println("Shutdown complete")
}

type App struct {

	// HTTP server for handling HTTP communication.
	HTTPServer *http.HttpServer

	CDSRepository *repository.CDSRepository
}

func NewApp() *App {
	httpServer := http.NewHttpServer()
	return &App{
		HTTPServer: httpServer,
	}
}

// RunHTTPServer executes the program. The configuration should already be set up before
// calling this function.
func (m *App) RunHTTPServer(ctx context.Context) (err error) {
	cfg := config.GetConfig()

	cdsRepository := repository.NewCDSRepository(cfg)

	// Setup HTTP server.
	// Attach services to main for testing.
	m.CDSRepository = &cdsRepository

	// Copy repository implementations to the HTTP server.
	m.HTTPServer.CDSRepository = cdsRepository

	m.HTTPServer.Router()
	logger.Info.Println("====* server ready *====")

	m.HTTPServer.Run()

	return nil
}

// Close gracefully stops the program.
func (m *App) Close() error {
	return nil
}
