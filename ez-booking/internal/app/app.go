package app

import (
	"context"
	"flag"
	"github.com/CyganFx/table-reservation/ez-booking/internal/cache"
	"github.com/CyganFx/table-reservation/ez-booking/internal/delivery/http-v1"
	"github.com/CyganFx/table-reservation/ez-booking/internal/repository/postgres"
	"github.com/CyganFx/table-reservation/ez-booking/internal/service"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/rest-errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Run initializes whole application
func Run(port string) {
	addr := flag.String("addr", port, "HTTP network address")
	dsn := flag.String("dsn",
		os.Getenv("POSTGRES_URI"),
		"PostgreSQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	dbPool, err := pgxpool.Connect(context.Background(), *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer dbPool.Close()

	templateCache, err := cache.NewTemplate("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	userRepo := postgres.NewUser(dbPool)
	reservationRepo := postgres.NewReservation(dbPool)
	userService := service.NewUser(userRepo)
	reservationService := service.NewReservation(reservationRepo)
	restErrorsResponser := rest_errors.NewHttpResponser(errorLog)
	handler := http_v1.NewHandler(userService, reservationService, restErrorsResponser, infoLog, templateCache)

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      handler.Init(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting http server on %v", *addr)
	errorLog.Fatal(srv.ListenAndServe())
}
