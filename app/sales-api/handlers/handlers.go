// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/ardanlabs/service/business/auth" // Import is removed in final PR
	"github.com/ardanlabs/service/business/mid"
	"github.com/ardanlabs/service/foundation/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, db *sqlx.DB, a *auth.Auth) http.Handler {
	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	// Register health check endpoint. This route is not authenticated.
	c := check{
		build: build,
		db:    db,
	}

	//Register static css wildcard
	app.Handle(http.MethodGet, "/css/*path", web.LoadCSS)

	app.Handle(http.MethodGet, "/v1/health", c.health)
	app.Handle(http.MethodGet, "/hello", hello)
	app.Handle(http.MethodGet, "/foo", foo)

	return app
}
