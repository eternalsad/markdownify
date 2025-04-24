package contract

import (
	"bytes"
	"github.com/eternalsad/markdownify/ast"
	"github.com/eternalsad/markdownify/md2"
	"github.com/eternalsad/markdownify/parser"
)

// Singleton instances for parser and renderer
var (
	mdParser    *parser.Parser
	md2Renderer *md2.Renderer
)

// init initializes the parser and renderer once
func init() {
	// Initialize parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	mdParser = parser.NewWithExtensions(extensions)

	// Initialize the renderer for Markdown V2
	md2Renderer = md2.NewRenderer()
}

// ConvertMD2 converts regular Markdown to Telegram's Markdown V2 format
func ConvertMD2(md string) string {
	// Parse the markdown input
	doc := mdParser.Parse([]byte(md))

	// Render to Markdown V2
	output := render(doc, md2Renderer)

	return string(output)
}

func render(doc ast.Node, renderer *md2.Renderer) []byte {
	var buf bytes.Buffer
	renderer.RenderHeader(&buf, doc)
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		return renderer.RenderNode(&buf, node, entering)
	})
	renderer.RenderFooter(&buf, doc)
	return buf.Bytes()
}
