package renderer

import "embed"

// Builtin contains the existing reference form, shared by both renderers.
//
//go:embed form.html form.js
var Builtin embed.FS
