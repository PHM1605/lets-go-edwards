package ui

import "embed"

// comment directive: tell Go to store files inside /static (ui/static) into "Files" variable

//go:embed "html" "static"
var Files embed.FS
