// Package swagger serves the raw OpenAPI spec and a Swagger UI page, and
// exposes a request-validation middleware built from the same embedded
// spec used to generate internal/generated.
package swagger

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/borkanie/brand-new-day-api/internal/generated"
)

// openapiYAML holds the raw contract bytes served at GET /openapi.yaml.
//
// This is a copy of the repo-root openapi.yaml kept in sync by `make
// generate` (see the Makefile's `generate` target, which copies
// openapi.yaml into this package directory right after codegen runs).
// go:embed cannot reach outside its own package directory, hence the copy.
//
//go:embed openapi.yaml
var openapiYAML []byte

// swaggerUIPage is a minimal HTML page that loads swagger-ui-dist from a CDN
// and points it at the /openapi.yaml route served by this package. No extra
// Go dependency is required for the UI itself.
const swaggerUIPage = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>brand_new_day API - Swagger UI</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: '/openapi.yaml',
        dom_id: '#swagger-ui',
      });
    };
  </script>
</body>
</html>
`

// MountRoutes registers the raw spec at GET /openapi.yaml and the Swagger UI
// at GET /swagger/ (and GET /swagger) on the supplied chi router.
func MountRoutes(router chi.Router) {
	router.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(openapiYAML)
	})

	serveSwaggerUIPage := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerUIPage))
	}

	router.Get("/swagger/", serveSwaggerUIPage)
	router.Get("/swagger", serveSwaggerUIPage)
}

// NewSpecValidatorMiddleware builds a chi-compatible middleware that
// validates incoming requests against the embedded OpenAPI spec (via
// generated.GetSwagger()) using kin-openapi/openapi3filter under the hood.
//
// The returned middleware wipes the spec's `Servers` field before building
// the validator so that Host-header validation (tied to the `servers:` entry
// in openapi.yaml) does not reject legitimate local requests.
func NewSpecValidatorMiddleware() (func(http.Handler) http.Handler, error) {
	swaggerSpec, err := generated.GetSwagger()
	if err != nil {
		return nil, err
	}

	return nethttpmiddleware.OapiRequestValidatorWithOptions(swaggerSpec, &nethttpmiddleware.Options{
		DoNotValidateServers: true,
	}), nil
}
