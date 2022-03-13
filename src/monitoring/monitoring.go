package monitoring

import (
	"time"

	"github.com/gin-gonic/gin"
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

var responsesSent = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "responses_sent",
	Help: "Count responses sent",
})

func Initialise() {
	prometheus.MustRegister(responsesSent)
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(cpuLoad)
	go func() {
		for {
			updateLoad()
			time.Sleep(time.Millisecond * 500)
		}
	}()
}

func UpdateResponseSent(c *gin.Context) {
	responsesSent.Add(1)

	c.Next()
}

func updateLoad() {
	cpuLoadTemp, _ := cpu.Percent(time.Second, false)
	cpuLoad.Set(cpuLoadTemp[0])
}
