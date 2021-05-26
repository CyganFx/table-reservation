package app

import (
	"github.com/CyganFx/table-reservation/internal/app/config"
	"github.com/CyganFx/table-reservation/internal/delivery/http-v1"
	"github.com/CyganFx/table-reservation/internal/repository/postgres"
	"github.com/CyganFx/table-reservation/internal/service"
	"github.com/CyganFx/table-reservation/pkg/cache"
	"github.com/CyganFx/table-reservation/pkg/rest-errors"
	"github.com/joho/godotenv"
	"log"
	"os"
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

	srv := new(config.Server)
	infoLog.Printf("main: API listening on host %s and port %s", cfg.Web.APIHost, cfg.Web.APIPort)
	if err := srv.Run(cfg, handler.Init(*cfg), errorLog); err != nil {
		errorLog.Fatal(err)
	}
}
