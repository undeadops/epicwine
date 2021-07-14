package main

import (
	"epicwine/pkg/metric"
	"os"
	"time"
)

// statrunner print stats
func statrunner(interrupt chan os.Signal, metric *metric.Metric) {
	// Ticker for every minute
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			metric.Output()
		case <-interrupt:
			break
		}
	}
}
