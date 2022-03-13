package monitoring

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

var cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius",
	Help: "Current temperature of the CPU.",
})

var cpuLoad = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_load_percentage",
	Help: "Current cpu load in percentage",
})

func Initialise() {
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(cpuLoad)
	go func() {
		for {
			cpuLoadTemp, _ := cpu.Percent(time.Second, false)
			cpuLoad.Set(cpuLoadTemp[0])
			time.Sleep(time.Millisecond * 500)
		}
	}()
}
