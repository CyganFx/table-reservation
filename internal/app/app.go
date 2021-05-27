package app

import (
	"context"
	"github.com/CyganFx/table-reservation/internal/app/config"
	"github.com/CyganFx/table-reservation/internal/delivery/http-v1"
	"github.com/CyganFx/table-reservation/internal/repository/postgres"
	"github.com/CyganFx/table-reservation/internal/service"
	"github.com/CyganFx/table-reservation/pkg/cache"
	"github.com/CyganFx/table-reservation/pkg/rest-errors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err.Error())
	}
}

// Run initializes whole application
func Run(configsDir, templatesDir string) {
	cfg, err := config.Init(configsDir)
	if err != nil {
		log.Fatalf("Error loading initializing configs: %v", err.Error())
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	templateCache, err := cache.NewTemplate(templatesDir)
	if err != nil {
		errorLog.Fatal(err)
	}
	dbPool, err := postgres.InitPool(*cfg)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer dbPool.Close()

	userRepo := postgres.NewUser(dbPool)
	reservationRepo := postgres.NewReservation(dbPool)
	userService := service.NewUser(userRepo)
	reservationService := service.NewReservation(reservationRepo)
	restErrorsResponser := rest_errors.NewHttpResponser(errorLog)
	handler := http_v1.NewHandler(userService, reservationService, restErrorsResponser, infoLog, templateCache)

	//Server
	srv := &http.Server{
		Addr:         ":" + cfg.Web.APIPort,
		ErrorLog:     errorLog,
		Handler:      handler.Init(*cfg),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	// Run Server
	go func() {
		infoLog.Printf("main: API listening on host %s and port %s", cfg.Web.APIHost, cfg.Web.APIPort)
		serverErrors <- srv.ListenAndServe()
	}()

	// Graceful Shutdown
	select {
	case err := <-serverErrors:
		errorLog.Fatalf("server error: %v", err)
	case sig := <-shutdown:
		infoLog.Printf("main: %v: Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			errorLog.Fatalf("could not stop server gracefully: %v", err)
		}

		infoLog.Printf("main: %v: Completed shutdown", sig)
	}
}
