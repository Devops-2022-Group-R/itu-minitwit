package controllers_test

import (
	"os"
	"testing"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
)

func TestMain(m *testing.M) {
	database.InitDatabase(database.TestDatabasePath)
	exitCode := m.Run()
	database.NukeDatabase(database.TestDatabasePath)
	os.Exit(exitCode)
}
