package utils

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func SanitizeInput(rawContent string) string {
	p := bluemonday.UGCPolicy()

	p.AllowElements("br")

	contentWithBreaks := strings.ReplaceAll(rawContent, "\n", "<br>")

	sanitizedHTML := p.Sanitize(contentWithBreaks)
	return sanitizedHTML
}
