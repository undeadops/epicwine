package server

import (
	"net/http"
	"time"
)

func (s *Server) totalmetrics(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Metric.AddTotal()
		h.ServeHTTP(w, r)
	})
}

func (s *Server) logging(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func(start time.Time) {
			requestID := w.Header().Get("X-Request-Id")
			if requestID == "" {
				requestID = "unknown"
			}
			s.Logger.Println(requestID, req.Method, req.URL.Path, req.RemoteAddr, req.UserAgent(), time.Since(start))
		}(time.Now())
		hdlr.ServeHTTP(w, req)
	})
}

func (s *Server) tracing(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := req.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = s.NextRequestID()
		}
		w.Header().Set("X-Request-Id", requestID)
		hdlr.ServeHTTP(w, req)
	})
}
