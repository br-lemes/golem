package assets

import "embed"

//go:embed css/output.css fonts/geist/*.woff2 images/* js/*.js
var Assets embed.FS
