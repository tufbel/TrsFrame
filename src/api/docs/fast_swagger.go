// Package docs
// Title       : fast_swagger.go
// Author      : Tuffy  2023/5/6 14:34
// Description :
package docs

import (
	"TrsFrame/src/tools/mylog"
	"fmt"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"html/template"
	"os"
	"regexp"
)

const fastHtml string = `
<!DOCTYPE html>
<html>

<head>
    <link type="text/css" rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui.css">
    <link rel="shortcut icon" href="https://fastapi.tiangolo.com/img/favicon.png">
    <title>NCC - Swagger UI</title>
</head>

<body>
    <div id="swagger-ui">
    </div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui-bundle.js"></script>
    <!-- "SwaggerUIBundle" is now available on the page -->
    <script>
    const ui = SwaggerUIBundle({
        url: '{{ .BaseURL }}/openapi.json',
        "dom_id": "#swagger-ui",
        "layout": "BaseLayout",
        "deepLinking": true,
        "showExtensions": true,
        "showCommonExtensions": true,
        oauth2RedirectUrl: window.location.origin + '{{ .BaseURL }}/docs/oauth2-redirect',
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
    })
    </script>
</body>

</html>
`

func OpenapiV3MarshalJSON(doc *openapi3.T) ([]byte, error) {
	m := make(map[string]interface{}, 4+len(doc.Extensions))
	for k, v := range doc.Extensions {
		m[k] = v
	}
	m["openapi"] = doc.OpenAPI
	if x := doc.Components; x != nil {
		m["components"] = x
	}
	m["info"] = doc.Info
	m["paths"] = doc.Paths
	if x := doc.Security; len(x) != 0 {
		m["security"] = x
	}
	if x := doc.Servers; len(x) != 0 {
		m["servers"] = x
	}
	if x := doc.Tags; len(x) != 0 {
		m["tags"] = x
	}
	if x := doc.ExternalDocs; x != nil {
		m["externalDocs"] = x
	}
	return json.MarshalIndent(m, "", "    ")
}

type FastSwagger struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description"`
	BaseURL     string `json:"base_url"`
	SwaggerPath string `json:"swagger_path"`
	OpenapiPath string `json:"openapi_path"`
	OpenapiJSON []byte `json:"openapi_json"`
}

func (f *FastSwagger) BuildOpenapi() {
	swaggerBytes, err := os.ReadFile(f.SwaggerPath)
	if err != nil {
		mylog.Error(fmt.Sprintf("Error reading %s file: %v", f.SwaggerPath, err))
		return
	}
	var swaggerDoc openapi2.T
	if err := json.Unmarshal(swaggerBytes, &swaggerDoc); err != nil {
		mylog.Error(fmt.Sprintf("Error unmarshaling SwaggerDoc: %v", err))
		return
	}

	openapiDoc, err := openapi2conv.ToV3(&swaggerDoc)
	if err != nil {
		mylog.Error(fmt.Sprintf("Error converting SwaggerDoc to OpenapiDoc: %v", err))
		return
	}
	// Set version
	openapiDoc.Info.Version = f.Version

	// Completion path
	newPaths := make(map[string]*openapi3.PathItem, len(openapiDoc.Paths))
	for path, pathItem := range openapiDoc.Paths {
		newPaths[f.BaseURL+path] = pathItem
	}
	openapiDoc.Paths = newPaths

	// Example Modify the BasePath
	domainRe := regexp.MustCompile(`^(https?://[^/]+)`)
	for _, server := range openapiDoc.Servers {
		server.URL = domainRe.FindString(server.URL)
	}

	openapiBytes, err := OpenapiV3MarshalJSON(openapiDoc)
	if err != nil {
		mylog.Error(fmt.Sprintf("Error marshaling OpenapiDoc: %v", err))
		return
	}
	if err := os.WriteFile(f.OpenapiPath, openapiBytes, 0644); err != nil {
		mylog.Error(fmt.Sprintf("Error writing %s file: %v", f.OpenapiPath, err))
		return
	}
}

func (f *FastSwagger) AddDocs(webRouter *gin.Engine) {
	indexTpl, _ := template.New("swagger_index.html").Parse(fastHtml)
	webRouter.GET(f.BaseURL+"/docs", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		indexTpl.Execute(ctx.Writer, f)
	})

	webRouter.GET(f.BaseURL+"/openapi.json", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		ctx.File(f.OpenapiPath)
	})
}
