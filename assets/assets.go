package assets

import (
	"embed"
)

//go:embed css images
var PublicAssetFS embed.FS

//go:embed templates
var TemplateFS embed.FS
