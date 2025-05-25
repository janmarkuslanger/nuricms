package embedfs

import "embed"

//go:embed templates/**/*.tmpl
var TemplatesFS embed.FS
