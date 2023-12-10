package assets

import (
	"embed"
)

//go:embed css
var PublicAssetFS embed.FS

//go:embed templates
var TemplateFS embed.FS
