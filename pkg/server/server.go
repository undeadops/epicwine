package server

import (
	"epicwine/pkg/csvfile"
	"epicwine/pkg/metric"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Metric        *metric.Metric
	Db            *csvfile.CsvFile
	Logger        *log.Logger
	NextRequestID func() string
}

func SetupServer(file string) *Server {
	return &Server{
		Metric:        metric.NewMetric(),
		Db:            csvfile.NewCsvFile(file),
		Logger:        log.New(os.Stdout, "", log.LstdFlags),
		NextRequestID: func() string { return strconv.FormatInt(time.Now().UnixNano(), 36) },
	}
}

// Router register necessary routes and returns an instance of a router.
func (s *Server) Router() *mux.Router {
	// Setup Mux Router
	r := mux.NewRouter()
	// Load up some http server middleware
	r.Use(s.tracing)
	r.Use(s.logging)
	r.Use(s.totalmetrics)

	r.HandleFunc("/status", s.status).Methods("GET")
	r.HandleFunc("/wine", s.listwine).Methods("GET")
	r.HandleFunc("/wine", s.putwine).Methods("PUT")
	r.HandleFunc("/wine/{id:[0-9]+}", s.getwine).Methods("GET")

	return r
}
