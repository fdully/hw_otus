package calendar

import "net/http"

func PreMuxRouter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.HandlerFunc(defaultCalendarHandler).ServeHTTP(w, r)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func defaultCalendarHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("welcome to calendar default page"))
}
