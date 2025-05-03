package rejoinder

import "github.com/andreas-bauer/rejoinderoo/internal/reader"

type rejoinder struct {
	file     reader.FileReader
	template string
}

func NewRejoinder() *rejoinder {
	// we need: selected, all headers and records, template

	return &rejoinder{}
}
