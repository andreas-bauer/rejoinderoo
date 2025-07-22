package project

import "embed"

var (
	//go:embed web/templates/*
	TemplateFS embed.FS
)
