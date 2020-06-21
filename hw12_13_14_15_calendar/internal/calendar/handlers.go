package calendar

import (
	"fmt"
	"net/http"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/util"
)

func (c Calendar) GetEventsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		events, err := c.r.GetEvents(r.Context())
		if err != nil {
			util.ServerError(r.Context(), w, err)
			return
		}
		_, _ = fmt.Fprint(w, events)
	})
}

func (c Calendar) DefaultHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("this is default calendar page"))
	})
}
