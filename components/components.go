package components

import (
	"github.com/gomarkdown/markdown/parser"
)

const extensions parser.Extensions = parser.CommonExtensions | parser.Mmark
