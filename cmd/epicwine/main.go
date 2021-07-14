package main

import (
	"context"
	"epicwine/pkg/server"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	port int
)

func main() {
	var f string

	// Parse Flags
	flag.IntVar(&port, "port", 5000, "Bind Server to Port")
	flag.StringVar(&f, "csv", "./wine.csv", "Location to the Wine CSV Datastore")
	flag.Parse()

	// Setup Connection timeouts
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Capture Sig Terms, send to sub goroutine
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Setup Server
	s := server.SetupServer(f)
	router := s.Router()

	// Verify Winelist Count
	wc, _ := s.Db.GetCount()
	s.Metric.SetWineCount(wc)

	// Configure HTTP Server
	s.Logger.Printf("Server is starting...")
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + strconv.Itoa(port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Stats Logger
	go statrunner(interrupt, s.Metric)

	// Startup HTTP Server
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	s.Logger.Printf("HTTP Server is ready")

	killSwitch := <-interrupt
	switch killSwitch {
	case os.Interrupt:
		s.Logger.Printf("Got SIGINT, Shutting down...")
	case syscall.SIGTERM:
		s.Logger.Printf("Got SIGTERM, Shutting down...")
	}
	srv.Shutdown(ctx)
	s.Logger.Printf("Shutdown Server")
}
