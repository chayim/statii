package plugins

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleConfig(t *testing.T) {
	sample := "statii.conf.example"
	here, _ := os.Getwd()
	samplefile := filepath.Join(here, sample)
	_, err := os.Stat(samplefile)
	if err != nil {
		samplefile = filepath.Join(here, "..", "..", sample)
	}
	cfg, _ := ioutil.ReadFile(samplefile)
	_, err = NewConfigFromBytes(cfg)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, err == nil)
}
