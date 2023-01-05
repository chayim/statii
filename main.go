package main

import (
	"flag"
	"os"
	"time"

	"github.com/chayim/statii/internal/plugins"
	"github.com/chayim/statii/internal/ui"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func main() {
	var config string
	flag.StringVar(&config, "c", "statii.conf", "path to configuration file")
	flag.Parse()

	_, err := os.Stat(config)
	if err != nil {
		log.Fatalf("%s does not exist", config)
	}

	data, err := os.ReadFile(config)
	if err != nil {
		log.Error(err)
		os.Exit(3)
	}
	cfg, err := plugins.NewConfigFromBytes(data)
	if err != nil {
		log.Error(err)
		os.Exit(3)
	}

	// handle scheduling, async, our plugin stores
	s := gocron.NewScheduler(time.UTC)
	s.Every(cfg.RescheduleEvery).Seconds().Do(cfg.ProcessPlugins)
	s.StartAsync()

	ui.RunApp(cfg)

}
