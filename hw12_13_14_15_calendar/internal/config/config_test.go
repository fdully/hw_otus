package config

import (
	"context"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {

	t.Run("init config", func(t *testing.T) {
		var configData = "{http_server: {host: localhost, port: 8080}, log_file: {path: /log/calendar.log, level: -1}," +
			"sql_db: {host: localhost, port: 5432, login: db_login, password: db_pass, database: test_db}," +
			"storage: {sql: true}, grpc_server: {host: localhost, port: 9111}}"

		sqlDB := SQLDB{
			Host:     "localhost",
			Port:     5432,
			Login:    "db_login",
			Password: "db_pass",
			Database: "test_db",
		}
		http := HTTPServer{
			Host: "localhost",
			Port: 8080,
		}
		grpc := GRPCServer{
			Host: "localhost",
			Port: 9111,
		}
		log := LogFile{
			Path:  "/log/calendar.log",
			Level: -1,
		}
		storage := Storage{
			SQL: true,
		}

		err := InitConfig(strings.NewReader(configData))
		require.NoError(t, err)

		ctx := context.Background()
		conf := FromContext(ctx)

		ctx = WithConfig(ctx, conf)

		require.Equal(t, sqlDB, FromContext(ctx).SQLDB)
		require.Equal(t, http, FromContext(ctx).HTTPServer)
		require.Equal(t, log, FromContext(ctx).LogFile)
		require.Equal(t, storage, FromContext(ctx).Storage)
		require.Equal(t, grpc, FromContext(ctx).GRPCServer)
	})
}
