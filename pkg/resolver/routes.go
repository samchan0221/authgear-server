package resolver

import (
	"net/http"

	"github.com/authgear/authgear-server/pkg/lib/config/configsource"
	"github.com/authgear/authgear-server/pkg/lib/deps"
	"github.com/authgear/authgear-server/pkg/resolver/handler"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/httputil"
)

func NewRouter(p *deps.RootProvider, configSource *configsource.ConfigSource) *httproute.Router {
	router := httproute.NewRouter()

	router.Add(httproute.Route{
		Methods:     []string{"GET"},
		PathPattern: "/healthz",
	}, http.HandlerFunc(httputil.HealthCheckHandler))

	chain := httproute.Chain(
		p.RootMiddleware(newPanicEndMiddleware),
		p.RootMiddleware(newPanicWriteEmptyResponseMiddleware),
		p.RootMiddleware(newBodyLimitMiddleware),
		p.RootMiddleware(newSentryMiddleware),
		&deps.RequestMiddleware{
			RootProvider: p,
			ConfigSource: configSource,
		},
		p.Middleware(newPanicLogMiddleware),
		p.Middleware(newSessionMiddleware),
	)

	route := httproute.Route{Middleware: chain}

	router.AddRoutes(p.Handler(newSessionResolveHandler), handler.ConfigureResolveRoute(route)...)

	return router
}
