package main

import (
	"flag"
	"os"

	"github.com/chayim/statii/internal/plugins"
	"github.com/chayim/statii/internal/ui"
	log "github.com/sirupsen/logrus"
)

func main() {
	var config string
	flag.StringVar(&config, "c", "statii.conf", "path to configuration file")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableLevelTruncation: true, ForceColors: true})
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}

	_, err := os.Stat(config)
	if err != nil {
		log.Fatalf("%s does not exist", config)
	}

	data, err := os.ReadFile(config)
	if err != nil {
		log.Error(err)
		os.Exit(3)
	}
	cfg, _ := plugins.NewConfigFromBytes(data)

	// handle scheduling, async, our plugin stores
	// s := gocron.NewScheduler(time.UTC)
	// s.Every(cfg.RescheduleEvery).Seconds().Do(cfg.ProcessPlugins)
	// s.StartAsync()

	ui.RunApp(cfg)

}
