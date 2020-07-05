package util

import (
	"context"
	"net/http"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
)

func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func ServerError(ctx context.Context, w http.ResponseWriter, err error) {
	logger := logging.FromContext(ctx)
	logger.Errorf("http server error", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func NotFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func NoContentYet(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
}
