package templates

import "embed"

//go:embed mail.html.gohtml
var EmailTemplates embed.FS

//go:embed mail.plain.gohtml
var EmailPlainTemplates embed.FS