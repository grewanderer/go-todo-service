package handlers

import (
	_ "embed"
	"net/http"
	"sync"

	"github.com/swaggest/swgui/v5cdn"
	"sigs.k8s.io/yaml"
)

//go:embed spec/openapi.yaml
var openapiYAML []byte

// DocsHandler serves the OpenAPI specification and Swagger UI.
type DocsHandler struct {
	ui       http.Handler
	jsonOnce sync.Once
	jsonData []byte
	jsonErr  error
}

// NewDocsHandler constructs documentation handlers.
func NewDocsHandler() *DocsHandler {
	return &DocsHandler{
		ui: v5cdn.New("go-todo-service API", "/docs/openapi.json", "/docs"),
	}
}

// UI serves the Swagger UI.
func (d *DocsHandler) UI(w http.ResponseWriter, r *http.Request) {
	d.ui.ServeHTTP(w, r)
}

// SpecYAML serves the raw OpenAPI YAML.
func (d *DocsHandler) SpecYAML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openapiYAML)
}

// SpecJSON serves the OpenAPI specification in JSON format.
func (d *DocsHandler) SpecJSON(w http.ResponseWriter, r *http.Request) {
	d.jsonOnce.Do(func() {
		d.jsonData, d.jsonErr = yaml.YAMLToJSON(openapiYAML)
	})
	if d.jsonErr != nil {
		respondError(w, r, http.StatusInternalServerError, "failed to render OpenAPI document")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(d.jsonData)
}
