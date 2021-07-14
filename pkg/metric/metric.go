package metric

import (
	"fmt"
	"os"
	"sync"
)

type Metric struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Error   int `json:"errors"`
	Wine    int `json:"num_wine"`
	mu      *sync.RWMutex
}

func NewMetric() *Metric {
	m := &Metric{
		mu: &sync.RWMutex{},
	}
	m.mu.Lock()
	m.setup()
	m.mu.Unlock()
	return m
}

// Output - Print out current metrics
// I Don't like this function
// TODO: Better way?
func (m *Metric) Output() {
	m.mu.Lock()
	// Generate Output For Logging Purposes
	av := (float64(m.Success) / float64(m.Total)) * 100
	fmt.Fprintf(os.Stderr, "REQUESTS=%v SUCCESS=%v ERRORS=%v AVAILABILITY=%v NUM_WINES=%v\n",
		m.Total, m.Success, m.Error, av, m.Wine)

	// created output, reset metrics
	m.setup()
	m.mu.Unlock()
}

func (m *Metric) SetWineCount(w int) {
	m.Wine = w
}

// AddTotal
func (m *Metric) AddTotal() {
	m.mu.Lock()
	m.Total = m.Total + 1
	m.mu.Unlock()
}

// Add Success
func (m *Metric) AddSuccess() {
	m.mu.Lock()
	m.Success = m.Success + 1
	m.mu.Unlock()
}

// Add Errors
func (m *Metric) AddErrors() {
	m.mu.Lock()
	m.Error = m.Error + 1
	m.mu.Unlock()
}

func (m *Metric) setup() {
	// reset stats to zero
	m.Total = 0
	m.Success = 0
	m.Error = 0
}
