// Package docs
// Title       : fast_openapi.go
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
	"github.com/swaggo/swag/gen"
	"html/template"
	"os"
	"path/filepath"
)

const fastHtml string = `
<!DOCTYPE html>
<html>

<head>
    <link type="text/css" rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui.css">
    <link rel="shortcut icon" href="https://fastapi.tiangolo.com/img/favicon.png">
    <title>{{ .Title }} - Docs</title>
</head>

<body>
    <div id="swagger-ui"></div>
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

type GenLogger struct{}

func (g GenLogger) Printf(format string, args ...interface{}) {
	mylog.Logger.Info(fmt.Sprintf(format, args...))
}

func BuildSwagger(apiDir string) error {
	return gen.New().Build(&gen.Config{
		SearchDir:           apiDir,
		Excludes:            "",
		ParseExtension:      "",
		MainAPIFile:         "main.go",
		PropNamingStrategy:  "camelcase",
		OutputDir:           filepath.Join(apiDir, "docs"),
		OutputTypes:         []string{"json"},
		ParseVendor:         false,
		ParseDependency:     false,
		MarkdownFilesDir:    "",
		ParseInternal:       false,
		GeneratedTime:       false,
		RequiredByDefault:   false,
		CodeExampleFilesDir: "",
		ParseDepth:          100,
		InstanceName:        "",
		OverridesFile:       ".swaggo",
		ParseGoList:         true,
		Tags:                "",
		LeftTemplateDelim:   "{{",
		RightTemplateDelim:  "}}",
		PackageName:         "",
		Debugger:            GenLogger{},
		CollectionFormat:    "csv",
	})
}

type FastOpenAPI struct {
	Title           string `json:"title"`
	Version         string `json:"version"`
	Description     string `json:"description"`
	BaseURL         string `json:"base_url"`
	ApiDir          string `json:"api_dir"`
	SwaggerFileName string `json:"swagger_file_name"`
	OpenapiFileName string `json:"openapi_file_name"`
}

func (f FastOpenAPI) BuildOpenapi() {
	if _, err := os.Stat(filepath.Join(f.ApiDir, "main.go")); err != nil {
		mylog.Info(fmt.Sprintf("'main.go' does not exist, do not build docs."))
		return
	}
	mylog.Info(fmt.Sprintf("Building docs..."))
	if err := BuildSwagger(f.ApiDir); err != nil {
		mylog.Error(fmt.Sprintf("Error building swagger: %v", err))
		return
	}

	swaggerPath := filepath.Join(f.ApiDir, "docs", f.SwaggerFileName)
	swaggerBytes, err := os.ReadFile(swaggerPath)
	if err != nil {
		mylog.Error(fmt.Sprintf("Error reading %s file: %v", swaggerPath, err))
		return
	}
	var swaggerDoc openapi2.T
	if err := json.Unmarshal(swaggerBytes, &swaggerDoc); err != nil {
		mylog.Error(fmt.Sprintf("Error unmarshaling SwaggerDoc: %v", err))
		return
	}

	mylog.Info(fmt.Sprintf("SwaggerDoc To OpenapiDoc..."))
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
	//domainRe := regexp.MustCompile(`^(https?://[^/]+)`)
	//for _, server := range openapiDoc.Servers {
	//	server.URL = domainRe.FindString(server.URL)
	//}
	openapiDoc.Servers = nil

	openapiBytes, err := OpenapiV3MarshalJSON(openapiDoc)
	if err != nil {
		mylog.Error(fmt.Sprintf("Error marshaling OpenapiDoc: %v", err))
		return
	}
	openapiPath := filepath.Join(f.ApiDir, "docs", f.OpenapiFileName)
	if err := os.WriteFile(openapiPath, openapiBytes, 0644); err != nil {
		mylog.Error(fmt.Sprintf("Error writing %s file: %v", openapiPath, err))
		return
	}
	mylog.Info(fmt.Sprintf("OpenapiDoc build success."))
}

func (f FastOpenAPI) AddDocs(webRouter *gin.Engine) {
	indexTpl, _ := template.New("swagger_index.html").Parse(fastHtml)
	webRouter.GET(f.BaseURL+"/docs", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		indexTpl.Execute(ctx.Writer, f)
	})

	openapiPath := filepath.Join(f.ApiDir, "docs", f.OpenapiFileName)
	webRouter.GET(f.BaseURL+"/openapi.json", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		ctx.File(openapiPath)
	})
}
