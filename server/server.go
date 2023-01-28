package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/casperprejler/node-exporter/metrics"
)

func SetupServer() {
	metricsRoute()

	http.ListenAndServe(":8080", nil)

}

type AllMetrics struct {
	CpuMetrics    metrics.CPUMetric    `json:"cpuMetrics"`
	MemoryMetrics metrics.MemoryMetric `json:"memMetrics"`
}

func metricsRoute() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {

		allMetrics := AllMetrics{
			CpuMetrics:    metrics.GetCPUMetrics(),
			MemoryMetrics: metrics.GetMemoryMetrics(),
		}
		metricsJson, err := json.Marshal(allMetrics)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(metricsJson))
	})
}
