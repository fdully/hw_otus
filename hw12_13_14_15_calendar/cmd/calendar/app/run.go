package app

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/api"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/api/swaggerui"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/pb"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/repository/memory"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/repository/sqldb"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/util/grpcutil"
	wu "github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/util/webutil"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/xlab/closer"
	"google.golang.org/grpc"
)

func RunCalendar(ctx context.Context) {
	logger := logging.FromContext(ctx)
	conf := config.FromContext(ctx)

	lsn, err := net.Listen("tcp", buildAddr(ctx, conf.GRPCServer.Host, conf.GRPCServer.Port))
	if err != nil {
		logger.Fatal("can't create grpc listener", err)
	}

	repo := defineRepository(ctx)
	cc := calendar.NewCalendar(repo)

	grpcGWwMux := runtime.NewServeMux()
	a := api.NewCalendarGRPCApi(cc)

	// registering grpc api
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcutil.UnaryInterceptor()))
	pb.RegisterCalendarServiceServer(grpcServer, a)

	// registering rest api
	err = pb.RegisterCalendarServiceHandlerServer(ctx, grpcGWwMux, a)
	if err != nil {
		logger.Fatal(err)
	}

	closer.Bind(func() {
		logger.Info("stopping grpc server")
		grpcServer.GracefulStop()
	})
	go func() {
		logger.Infof("starting grpc server on port %d", conf.GRPCServer.Port)
		err := grpcServer.Serve(lsn)
		if err != nil {
			logger.Fatal("grpc serve", err)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/", wu.LogHTTPRequests(calendar.PreMuxRouter(grpcGWwMux)))
	mux.Handle("/swaggerui/", wu.LogHTTPRequests(http.StripPrefix("/swaggerui/", swaggerui.NewHandler(ctx))))

	webServer := wu.NewWebServer(mux, buildAddr(ctx, conf.HTTPServer.Host, conf.HTTPServer.Port))
	closer.Bind(func() {
		logger.Info("stopping web server")
		if err := webServer.Shutdown(time.Second * 2); err != nil {
			logger.Error("shutdown web server", err)
		}
	})
	logger.Infof("starting http server on port %d", conf.HTTPServer.Port)
	if err := webServer.Start(); err != nil {
		logger.Fatal("run calendar", err)
	}

	closer.Hold()
}

func defineRepository(ctx context.Context) calendar.Repository {
	logger := logging.FromContext(ctx)
	conf := config.FromContext(ctx)

	var repo calendar.Repository

	if conf.Storage.SQL {
		sqlPool, err := sqldb.OpenDB(ctx)
		if err != nil {
			logger.Fatal("connecting to sql db", err)
		}
		closer.Bind(func() {
			logger.Info("stopping sql pool")
			sqlPool.Close()
		})

		repo = sqldb.Repo{Pool: sqlPool}
	} else {
		repo = memory.NewRepo()
	}

	return repo
}

func buildAddr(ctx context.Context, host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}
