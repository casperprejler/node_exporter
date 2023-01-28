package metrics

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CPUMetric struct {
	CPUTotal       float64 `json:"cpuTotal"`
	CPUStateUser   float64 `json:"cpuStateUser"`
	CPUStateSystem float64 `json:"cpuStateSystem"`
	CPUStateIdle   float64 `json:"cpuStateIdle"`
	CPUStateIOWait float64 `json:"cpuStateIOWait"`
}

type MemoryMetric struct {
	MemTotal     int64 `json:"memTotal"`
	MemFree      int64 `json:"memFree"`
	MemAvailable int64 `json:"memAvailable"`
	MemCached    int64 `json:"memCached"`
	SwapTotal    int64 `json:"swapTotal"`
	SwapFree     int64 `json:"swapFree"`
}

func GetCPUMetrics() CPUMetric {

	// Read the contents of the /proc/stat file
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	line := scanner.Text()

	fields := strings.Fields(line)
	// Get CPU usage in jiffies
	user, _ := strconv.ParseUint(fields[1], 10, 64)
	nice, _ := strconv.ParseUint(fields[2], 10, 64)
	system, _ := strconv.ParseUint(fields[3], 10, 64)
	idle, _ := strconv.ParseUint(fields[4], 10, 64)
	iowait, _ := strconv.ParseUint(fields[5], 10, 64)
	irq, _ := strconv.ParseUint(fields[6], 10, 64)
	softirq, _ := strconv.ParseUint(fields[7], 10, 64)
	total := user + nice + system + idle + iowait + irq + softirq

	// Get CPU ClockRate
	clockRate := getCPUClockRate()

	userSeconds := float64(user) / float64(clockRate)
	systemSeconds := float64(system) / float64(clockRate)
	idleSeconds := float64(idle) / float64(clockRate)
	ioWaitSeconds := float64(iowait) / float64(clockRate)
	totalSeconds := float64(total) / float64(clockRate)

	cpuMetric := CPUMetric{
		CPUTotal:       totalSeconds,
		CPUStateUser:   userSeconds,
		CPUStateSystem: systemSeconds,
		CPUStateIdle:   idleSeconds,
		CPUStateIOWait: ioWaitSeconds,
	}
	return cpuMetric
}

func GetMemoryMetrics() MemoryMetric {

	memoryMetric := MemoryMetric{
		MemTotal:     getIndividualMemoryBytes("MemTotal"),
		MemFree:      getIndividualMemoryBytes("MemFree"),
		MemAvailable: getIndividualMemoryBytes("MemAvailable"),
		MemCached:    getIndividualMemoryBytes("Cached"),
		SwapTotal:    getIndividualMemoryBytes("SwapTotal"),
		SwapFree:     getIndividualMemoryBytes("SwapFree"),
	}

	fmt.Print(memoryMetric)
	return memoryMetric

}

func getIndividualMemoryBytes(metric string) int64 {
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
	return memoryInBytes

}

func getCPUClockRate() float64 {
	var clockRate float64
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu MHz") {
			fields := strings.Split(line, ":")
			if len(fields) == 2 {
				clockRateString := strings.TrimSpace(fields[1])
				clockRate, err = strconv.ParseFloat(clockRateString, 64)
				if err != nil {
					log.Println(err)
				}
				break
			}
		}
	}
	return clockRate
}
