package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/cmd/calendar/app"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "-config /configs/calendar.cfg")

	flag.Parse()

	f, err := os.Open(configFile)
	if err != nil {
		flag.Usage()
		log.Fatalf("ERROR: can't open cfgfile %v\n", err)
	}

	ctx := context.Background()

	if err := config.InitConfig(f); err != nil {
		log.Fatalf("ERROR: can't init config %v\n", err)
	}
	conf := config.FromContext(ctx)

	if err := logging.InitLog(conf.LogFile.Level, conf.LogFile.Path); err != nil {
		log.Fatalf("ERROR: can't init logging %v\n", err)
	}
	logger := logging.FromContext(ctx)

	ctx = logging.WithLogger(ctx, logger)
	ctx = config.WithConfig(ctx, conf)

	app.RunCalendar(ctx)
}
