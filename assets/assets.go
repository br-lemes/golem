package assets

import "embed"

//go:embed golem.css golem.js images/* tabler-icons-* theme.js
var Assets embed.FS
