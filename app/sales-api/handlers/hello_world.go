package handlers

import (
	"context"
	"net/http"

	"github.com/ardanlabs/service/foundation/web"
	"go.opentelemetry.io/otel/api/global"
)

func hello(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := global.Tracer("service").Start(ctx, "handlers.check.health")
	defer span.End()

	return web.Respond(ctx, w, "world", http.StatusOK)
}
