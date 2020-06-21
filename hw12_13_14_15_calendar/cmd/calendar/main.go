package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/repository/memory"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/repository/sqldb"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/util"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/webserver"

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

	if err := config.InitConfig(f); err != nil {
		log.Fatalf("ERROR: can't init config %v\n", err)
	}

	ctx := context.Background()
	conf := config.FromContext(ctx)

	if err := logging.InitLog(conf.LogFile.Level, conf.LogFile.Path); err != nil {
		log.Fatalf("ERROR: can't init logging %v\n", err)
	}
	logger := logging.FromContext(ctx)

	var repo calendar.Repository

	if conf.Storage.SQL {
		sqlPool, err := sqldb.OpenDB(ctx, fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable", conf.SQLDB.Host, conf.SQLDB.Port, conf.SQLDB.Login, conf.SQLDB.Password, conf.SQLDB.Database))
		if err != nil {
			logger.Fatal("connecting to sql db", err)
		}
		repo = sqldb.Repo{Pool: sqlPool}
	} else if conf.Storage.Memory {
		repo = memory.NewRepo()
	}

	c, err := calendar.NewCalendar(repo)
	if err != nil {
		logger.Fatalf("can't create calendar instance", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", c.DefaultHandler())
	mux.Handle("/events/", util.LogHTTPRequests(c.GetEventsHandler()))

	webServer := webserver.NewWebServer(mux, conf.HTTPServer.Address+":"+strconv.Itoa(conf.HTTPServer.Port))
	if err := webServer.Start(); err != nil {
		logger.Fatalf("problem with web server", err)
	}
}
