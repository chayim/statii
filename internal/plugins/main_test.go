package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	here, _ := os.Getwd()
	godotenv.Load(".env")
	godotenv.Load(filepath.Join(here, "..", "..", ".env"))
	m.Run()
}
