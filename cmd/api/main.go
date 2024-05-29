package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

const defaultIdleTimeout = time.Minute

const defaultReadTimeout = 10 * time.Second
const defaultWriteTimeout = 30 * time.Second

type config struct {
	port         int
	env          string
	idleTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type application struct {
	config config
	logger *log.Logger
}

func parseFlags(cfg *config) {
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.DurationVar(&cfg.idleTimeout, "idle-timeout", defaultIdleTimeout, "Server idle timeout")
	flag.DurationVar(&cfg.readTimeout, "read-timeout", defaultReadTimeout, "Server read timeout")
	flag.DurationVar(&cfg.writeTimeout, "write-timeout", defaultWriteTimeout, "Server write timeout")
	flag.Parse()
}

func main() {
	var cfg config

	parseFlags(&cfg)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  cfg.idleTimeout,
		ReadTimeout:  cfg.readTimeout,
		WriteTimeout: cfg.writeTimeout,
	}

	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
