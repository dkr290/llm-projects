package helpers

import (
	"log/slog"
	"net/http"
)

func MakeHandler(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {

			slog.Error("internal server error", "err", err, "path", r.URL.Path)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
