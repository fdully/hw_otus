package swaggerui

import (
	"context"
	"net/http"

	// init static files for swagger ui
	_ "github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/api/swaggerui/statik"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
	"github.com/rakyll/statik/fs"
)

// openapi ui handler on /swaggerui.
func NewHandler(ctx context.Context) http.Handler {
	logger := logging.FromContext(ctx)

	statikFS, err := fs.New()
	if err != nil {
		logger.Fatalf("can't do statik openapiui files binding", err)
	}

	return http.FileServer(statikFS)
}
