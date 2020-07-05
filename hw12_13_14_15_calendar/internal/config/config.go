package config

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.uber.org/config"
)

type Configuration struct {
	LogFile    LogFile
	SQLDB      SQLDB
	HTTPServer HTTPServer
	Storage    Storage
}

type Storage struct {
	SQL    bool
	Memory bool
}

type SQLDB struct {
	Host     string
	Port     int
	Login    string
	Password string
	Database string
}

type LogFile struct {
	Path  string
	Level int
}

type HTTPServer struct {
	Address string
	Port    int
}

type configKey struct{}

var c Configuration

func InitConfig(r io.Reader) error {
	provider, err := config.NewYAML(config.Expand(os.LookupEnv), config.Source(r))
	if err != nil {
		return err
	}

	if err := provider.Get("storage").Populate(&c.Storage); err != nil {
		return fmt.Errorf("ERROR: parsing storage config %w", err)
	}
	if (c.Storage.SQL && c.Storage.Memory) || (!c.Storage.Memory && !c.Storage.SQL) {
		return fmt.Errorf("please define correct storage. It must be either memory or sql. It can't be both true or both false")
	}

	if err := provider.Get("log_file").Populate(&c.LogFile); err != nil {
		return fmt.Errorf("ERROR: parsing log_file config %w", err)
	}

	if c.Storage.SQL {
		if err := provider.Get("sql_db").Populate(&c.SQLDB); err != nil {
			return fmt.Errorf("ERROR: parsing sql_db config %w", err)
		}
	}

	if err := provider.Get("http_server").Populate(&c.HTTPServer); err != nil {
		return fmt.Errorf("ERROR: parsing http_server config %w", err)
	}

	return nil
}

func WithConfig(ctx context.Context, c *Configuration) context.Context {
	return context.WithValue(ctx, configKey{}, c)
}

func FromContext(ctx context.Context) *Configuration {
	if conf, ok := ctx.Value(configKey{}).(*Configuration); ok {
		return conf
	}
	return &c
}
