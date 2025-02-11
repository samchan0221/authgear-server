package middleware

import (
	"net/http"

	"github.com/iawaknahc/originmatcher"

	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/util/log"
)

type CORSMiddlewareLogger struct{ *log.Logger }

func NewCORSMiddlewareLogger(lf *log.Factory) CORSMiddlewareLogger {
	return CORSMiddlewareLogger{lf.New("cors-middleware")}
}

type CORSMiddleware struct {
	Config *config.HTTPConfig
	Logger CORSMiddlewareLogger
}

func (m *CORSMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matcher, err := originmatcher.New(m.Config.AllowedOrigins)
		// nolint: staticcheck
		if err != nil {
			// FIXME(logging): Log invalid AllowedOrigins error here.
		}

		w.Header().Add("Vary", "Origin")

		origin := r.Header.Get("Origin")
		if origin != "" && err == nil && matcher.MatchOrigin(origin) {
			corsMethod := r.Header.Get("Access-Control-Request-Method")
			corsHeaders := r.Header.Get("Access-Control-Request-Headers")

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "900") // 15 mins

			if corsMethod != "" {
				w.Header().Set("Access-Control-Allow-Methods", corsMethod)
			}

			if corsHeaders != "" {
				w.Header().Set("Access-Control-Allow-Headers", corsHeaders)
			}
		}

		requestMethod := r.Method
		if requestMethod == "OPTIONS" {
			// FIXME(logging): the log below may cause too many logs to server
			// should remove the log after diagnostic
			m.Logger.Info("return 200 for options request")
			w.WriteHeader(http.StatusOK)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
