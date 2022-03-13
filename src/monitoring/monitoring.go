package monitoring

import (
	"log"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/prometheus/client_golang/prometheus"
)

var cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius",
	Help: "Current temperature of the CPU.",
})

var UserCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "minitwit_user_count",
	Help: "The amount of registered MiniTwit users",
})

func Initialise(openDatabase database.OpenDatabaseFunc) {
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(UserCount)
	initUserCount(openDatabase)
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
