package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/chayim/statii/internal/comms"
	"github.com/chayim/statii/internal/plugins"
	"github.com/chayim/statii/internal/ui"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func main() {
	var config string
	var fill bool
	var uiOnly bool
	var dataOnly bool
	var validateCfg bool
	var wipeDB bool
	flag.StringVar(&config, "c", "statii.conf", "path to configuration file")
	flag.BoolVar(&validateCfg, "v", false, "validate the configuration file")
	flag.BoolVar(&fill, "x", false, "fill the database junk sample data - and exit")
	flag.BoolVar(&dataOnly, "d", false, "process plugins only (ie no UI)")
	flag.BoolVar(&uiOnly, "u", false, "display the UI only (ie no data insertions)")
	flag.BoolVar(&wipeDB, "w", false, "wipe the database and exit")
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
	cfg, err := plugins.NewConfigFromBytes(data)
	if validateCfg {
		if err == nil {
			log.Info("configuration is valid, exiting.")
			os.Exit(0)
		}
	}
	if wipeDB {
		con := comms.NewConnection(cfg.Database, cfg.Size)
		_, err := con.DB.FlushAll(context.TODO()).Result()
		if err != nil {
			log.Error("failed to flush database")
			os.Exit(1)
		}
		os.Exit(0)
	}

	if fill {
		con := comms.NewConnection(cfg.Database, cfg.Size)
		var messages []*comms.Message
		for i := 0; i <= 25; i++ {
			m := comms.NewMessage("1",
				"fakedata",
				fmt.Sprintf("I am title %d", i),
				fmt.Sprintf("http://google.ca/search?q=%d", i),
				"green",
			)
			messages = append(messages, m)
		}
		con.SaveMany(context.TODO(), messages)
		os.Exit(0)
	}

	// when dataonly is true, just process
	// if both are faulse, they're default so run through
	if dataOnly || (!uiOnly && !dataOnly) {
		// handle scheduling, async, our plugin stores
		s := gocron.NewScheduler(time.UTC)
		s.Every(cfg.RescheduleEvery).Seconds().Do(cfg.ProcessPlugins)
		s.StartAsync()
	}

	// when uionly is true, just display
	// when both are false, they're default and not set - so run
	if uiOnly || (!uiOnly && !dataOnly) {
		ui.RunApp(cfg)
	}

}
