package config

import (
	"context"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {

	t.Run("init config", func(t *testing.T) {
		var configData = "{http_server: {address: localhost, port: 8080}, log_file: {path: /log/calendar.log, level: -1}," +
			"sql_db: {host: localhost, port: 5432, login: db_login, password: db_pass, database: test_db}," +
			"storage: {sql: true, memory: false}}"

		sqlDB := SQLDB{
			Host:     "localhost",
			Port:     5432,
			Login:    "db_login",
			Password: "db_pass",
			Database: "test_db",
		}
		http := HTTPServer{
			Address: "localhost",
			Port:    8080,
		}
		log := LogFile{
			Path:  "/log/calendar.log",
			Level: -1,
		}
		storage := Storage{
			SQL:    true,
			Memory: false,
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
	})

	t.Run("only memory storage", func(t *testing.T) {
		var configData = "{http_server: {address: localhost, port: 8080}, log_file: {path: /log/calendar.log, level: -1}," +
			"sql_db: {host: localhost, port: 5432, login: db_login, password: db_pass, database: test_db}," +
			"storage: {sql: false, memory: true}}"

		http := HTTPServer{
			Address: "localhost",
			Port:    8080,
		}
		log := LogFile{
			Path:  "/log/calendar.log",
			Level: -1,
		}
		storage := Storage{
			SQL:    false,
			Memory: true,
		}

		err := InitConfig(strings.NewReader(configData))
		require.NoError(t, err)

		ctx := context.Background()
		conf := FromContext(ctx)

		ctx = WithConfig(ctx, conf)

		require.Equal(t, http, FromContext(ctx).HTTPServer)
		require.Equal(t, log, FromContext(ctx).LogFile)
		require.Equal(t, storage, FromContext(ctx).Storage)
	})

	t.Run("both storages false or both true", func(t *testing.T) {
		configData := "{http_server: {address: localhost, port: 8080}, log_file: {path: /log/calendar.log, level: -1}," +
			"sql_db: {host: localhost, port: 5432, login: db_login, password: db_pass, database: test_db}," +
			"storage: {sql: false, memory: false}}"

		err := InitConfig(strings.NewReader(configData))
		require.Error(t, err)

		configData = "{http_server: {address: localhost, port: 8080}, log_file: {path: /log/calendar.log, level: -1}," +
			"sql_db: {host: localhost, port: 5432, login: db_login, password: db_pass, database: test_db}," +
			"storage: {sql: true, memory: true}}"

		err = InitConfig(strings.NewReader(configData))
		require.Error(t, err)
	})
}
