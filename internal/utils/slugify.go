package utils

import (
	"regexp"
	"strings"
)

func Slugify(text string) string {
	slug := strings.ReplaceAll(text, " ", "-")
	slug = strings.ToLower(slug)
	reg := regexp.MustCompile(`[^\w\-]+`)
	slug = reg.ReplaceAllString(slug, "")
	slug = strings.Trim(slug, "-")

	return slug
}
