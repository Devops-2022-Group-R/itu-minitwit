package monitoring

import (
	"time"

	"log"

	"github.com/gin-gonic/gin"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

var cpuLoad = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
	Name: "minitwit_cpu_load_percentage",
	Help: "Current cpu load in percentage",
}, getCpuLoad)

var responsesSent = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "minitwit_responses_sent",
	Help: "Count responses sent",
})

var requestDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "minitwit_request_duration",
	Help:    "Request duration in Milliseconds",
	Buckets: []float64{50.0, 100.0, 200.0, 500.0, 1000.0},
})

var UserCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "minitwit_user_count",
	Help: "The amount of registered MiniTwit users",
})

func Initialise(openDatabase database.OpenDatabaseFunc) {
	prometheus.MustRegister(responsesSent)
	prometheus.MustRegister(requestDurationHistogram)
	prometheus.MustRegister(cpuLoad)
	prometheus.MustRegister(UserCount)
	initUserCount(openDatabase)
}

func UpdateResponseSent(c *gin.Context) {
	responsesSent.Inc()
	c.Next()
}

func getCpuLoad() float64 {
	cpuLoadTemp, _ := cpu.Percent(time.Second, false)
	return cpuLoadTemp[0]
}

func RequestDuration(c *gin.Context) {
	startTime := time.Now()

	c.Next()

	duration := time.Since(startTime).Milliseconds()
	requestDurationHistogram.Observe(float64(duration))
}

func initUserCount(openDatabase database.OpenDatabaseFunc) {
	gormDb, err := database.ConnectDatabase(openDatabase)
	if err != nil {
		log.Fatal(err)
	}

	userRepository := database.NewGormUserRepository(gormDb)

	numUsers, err := userRepository.NumUsers()
	if err != nil {
		log.Fatal(err)
	}

	UserCount.Set(float64(numUsers))
}
