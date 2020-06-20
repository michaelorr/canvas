package web

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/spf13/afero"
	"go.opentelemetry.io/otel/api/global"
)

const assetPath = `/app/static`

func Render(path string) ([]byte, error) {
	var markup bytes.Buffer

	tpl, err := loadTemplate(path)
	if err != nil {
		return nil, fmt.Errorf("error processing template : %v", err)
	}

	if err = tpl.Execute(&markup, nil); err != nil {
		return nil, fmt.Errorf("error processing template : %v", err)
	}

	return markup.Bytes(), nil
}

func loadTemplate(path string) (*template.Template, error) {
	tmpl := template.Must(template.ParseFiles(filepath.Join(assetPath, path)))
	return tmpl, tmpl.Execute(os.Stdout, nil)
}

func LoadCSS(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := global.Tracer("service").Start(ctx, "handlers.check.health")
	defer span.End()

	var css bytes.Buffer

	fs := afero.NewBasePathFs(afero.NewOsFs(), assetPath)
	path := filepath.Join("css", httptreemux.ContextParams(r.Context())["path"])

	exists, err := afero.Exists(fs, path)
	if !exists {
		return NewRequestError(err, http.StatusNotFound)
	}
	if err != nil {
		return err
	}

	file, err := fs.Open(path)
	if err != nil {
		return err
	}

	_, err = css.ReadFrom(file)
	if err != nil {
		return err
	}

	v, _ := ctx.Value(KeyValues).(*Values)
	v.StatusCode = http.StatusOK
	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)

	_, err = css.WriteTo(w)
	return err
}
