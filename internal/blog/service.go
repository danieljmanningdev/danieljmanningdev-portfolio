// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package blog

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/yuin/goldmark"
)

func RenderMarkdown(content string) (template.HTML, error) {
	var buf bytes.Buffer

	if err := goldmark.Convert(
		[]byte(content),
		&buf,
	); err != nil {
		return "",
			fmt.Errorf("render markdown: %w", err)
	}

	return template.HTML(buf.String()), nil
}
