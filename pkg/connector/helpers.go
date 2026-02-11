package connector

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const ResourcesPageSize = 50

func titleCase(s string) string {
	titleCaser := cases.Title(language.English)

	return titleCaser.String(s)
}
