package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/repository/memory"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/repository/sqldb"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/sethvargo/go-signalcontext"
)

func main() {
	ctx, done := signalcontext.OnInterrupt()
	defer done()

	var configFile string
	flag.StringVar(&configFile, "config", "", "-config /configs/calendar.cfg")

	flag.Parse()

	f, err := os.Open(configFile)
	if err != nil {
		flag.Usage()
		log.Fatalf("ERROR: can't open cfgfile %v\n", err)
	}

	if err := config.InitConfig(f); err != nil {
		log.Fatalf("ERROR: can't init config %v\n", err)
	}
	conf := config.FromContext(ctx)

	if err := logging.InitLog(conf.LogFile.Level, conf.LogFile.Path); err != nil {
		log.Fatalf("ERROR: can't init logging %v\n", err)
	}
	logger := logging.FromContext(ctx)

	logger.Info("starting application")

	ctx = logging.WithLogger(ctx, logger)
	ctx = config.WithConfig(ctx, conf)

	if err := realMain(ctx); err != nil {
		logger.Fatal(err)
	}
}

func realMain(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	conf := config.FromContext(ctx)

	var repo calendar.Repository

	if conf.Storage.SQL {
		sqlPool, err := sqldb.OpenDB(ctx)
		if err != nil {
			logger.Fatal("connecting to sql db", err)
		}
		defer func() {
			logger.Info("stopping sql pool")
			sqlPool.Close()
		}()

		repo = sqldb.Repo{Pool: sqlPool}
	} else {
		repo = memory.NewRepo()
	}

	q := rabbit.NewConnector(ctx, conf.AMQP.URL, conf.AMQP.Exchange, conf.AMQP.Name,
		conf.AMQP.QOS)
	defer q.Close()

	s := scheduler.NewScheduler(repo, q)

	return s.Run(ctx)
}
