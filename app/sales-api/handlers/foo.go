package handlers

import (
	"context"
	"net/http"

	"github.com/ardanlabs/service/foundation/web"
	"go.opentelemetry.io/otel/api/global"
)

func foo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := global.Tracer("service").Start(ctx, "handlers.check.health")
	defer span.End()

	markup, err := web.Render(`foo.html`)
	if err != nil {
		return err
	}

	web.RespondHTML(ctx, w, markup, http.StatusOK)
	return nil
}
