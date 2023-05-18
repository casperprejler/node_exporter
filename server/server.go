package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/casperprejler/node-exporter/metrics"
)

func SetupServer() {
	metricsRoute()
	log.Print("Serving on :8080")
	http.ListenAndServe(":8080", nil)

}

type CollectedMetrics struct {
	CpuMetrics    metrics.CPUMetric
	MemoryMetrics metrics.MemoryMetric
}

func (c CollectedMetrics) String() string {
	return fmt.Sprintf("%s%s",
		c.CpuMetrics.String(),
		c.MemoryMetrics.String())
}

func metricsRoute() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {

		cMetrics := CollectedMetrics{
			CpuMetrics:    metrics.GetCPUMetrics(),
			MemoryMetrics: metrics.GetMemoryMetrics(),
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(cMetrics.String()))
	})
}
