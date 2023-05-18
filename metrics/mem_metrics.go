package metrics

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type MemoryMetric struct {
	MemTotal     string
	MemFree      string
	MemAvailable string
	MemCached    string
	SwapTotal    string
	SwapFree     string
}

func (m MemoryMetric) String() string {
	return fmt.Sprintf("mem_total %s\n"+
		"mem_bytes_free %s\n"+
		"mem_bytes_available %s\n"+
		"mem_bytes_cached %s\n"+
		"mem_swap_total %s\n"+
		"mem_swap_free %s\n",
		m.MemTotal,
		m.MemFree,
		m.MemAvailable,
		m.MemCached,
		m.MemFree,
		m.SwapTotal)
}

func GetMemoryMetrics() MemoryMetric {

	memoryMetric := MemoryMetric{
		MemTotal:     getMemMetricInBytes("MemTotal"),
		MemFree:      getMemMetricInBytes("MemFree"),
		MemAvailable: getMemMetricInBytes("MemAvailable"),
		MemCached:    getMemMetricInBytes("Cached"),
		SwapTotal:    getMemMetricInBytes("SwapTotal"),
		SwapFree:     getMemMetricInBytes("SwapFree"),
	}

	return memoryMetric

}

func getMemMetricInBytes(metric string) string {
	var memoryInBytes int64
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, metric) {
			fields := strings.Split(line, ":")
			if len(fields) == 2 {
				firstTrim := strings.TrimSpace(fields[1])
				secondSplit := strings.Split(firstTrim, " ")
				valueStr := strings.TrimSpace(secondSplit[0])

				value, err := strconv.ParseInt(valueStr, 10, 64)
				memoryInBytes = value / 1024
				if err != nil {
					log.Print(err)
				}
				break
			}
		}
	}
	return strconv.Itoa(int(memoryInBytes))

}
