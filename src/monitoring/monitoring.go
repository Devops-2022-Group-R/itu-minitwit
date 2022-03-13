package monitoring

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

var cpuLoad = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "minitwit_cpu_load_percentage",
	Help: "Current cpu load in percentage",
})

var responsesSent = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "minitwit_responses_sent",
	Help: "Count responses sent",
})

var requestDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "minitwit_request_duration",
	Help:    "Request duration in Milliseconds",
	Buckets: []float64{.05, .1, .2, .3, .5, 1.0, 2.0, 3},
})

func Initialise() {
	prometheus.MustRegister(responsesSent)
	prometheus.MustRegister(requestDurationHistogram)
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

func RequestDuration(c *gin.Context) {
	startTime := time.Now()

	c.Next()

	duration := time.Since(startTime)
	requestDurationHistogram.Observe(float64(duration))
}
