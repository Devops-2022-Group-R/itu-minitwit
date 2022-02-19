package controllers_test

import (
	"os"
	"testing"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openDatabase() gorm.Dialector {
	return sqlite.Open(":memory:")
}

func TestMain(m *testing.M) {
	database.InitDatabase(openDatabase)
	exitCode := m.Run()
	os.Exit(exitCode)
}
