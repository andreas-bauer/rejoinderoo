package project

import "embed"

var (
	//go:embed web/templates/*
	TemplateFS embed.FS

	//go:embed web/css/output.css
	CSS embed.FS
)
