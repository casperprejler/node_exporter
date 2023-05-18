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
	CPUTotal       string
	CPUStateUser   string
	CPUStateSystem string
	CPUStateIdle   string
	CPUStateIOWait string
}

func (m CPUMetric) String() string {
	return fmt.Sprintf("cpu_total %s\n"+
		"cpu_state_user %s\n"+
		"cpu_state_system %s\n"+
		"cpu_state_idle %s\n"+
		"cpu_state_iowait %s\n",
		m.CPUTotal,
		m.CPUStateUser,
		m.CPUStateSystem,
		m.CPUStateIdle,
		m.CPUStateIOWait)
}

// On UNIX systems a jiffie is a measurement to represent the amount of time the CPU
// has spent in terms of clock ticks.

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

	userSeconds := strconv.FormatFloat(float64(user)/float64(clockRate), 'f', -1, 64)
	systemSeconds := strconv.FormatFloat(float64(system)/float64(clockRate), 'f', -1, 64)
	idleSeconds := strconv.FormatFloat(float64(idle)/float64(clockRate), 'f', -1, 64)
	ioWaitSeconds := strconv.FormatFloat(float64(iowait)/float64(clockRate), 'f', -1, 64)
	totalSeconds := strconv.FormatFloat(float64(total)/float64(clockRate), 'f', -1, 64)

	cpuMetric := CPUMetric{
		CPUTotal:       totalSeconds,
		CPUStateUser:   userSeconds,
		CPUStateSystem: systemSeconds,
		CPUStateIdle:   idleSeconds,
		CPUStateIOWait: ioWaitSeconds,
	}

	return cpuMetric

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
