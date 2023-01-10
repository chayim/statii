package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

// ensures our dotnev file is loaded.. so that testing secrets happens outside of the repository
func TestMain(m *testing.M) {
	here, _ := os.Getwd()
	godotenv.Load(".env")
	godotenv.Load(filepath.Join(here, "..", "..", ".env"))
	m.Run()
}
