package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	mathjax "github.com/litao91/goldmark-mathjax"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
    chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
)


type MarkdownConverter struct{
	md goldmark.Markdown
}

func NewMarkdownConverter() *MarkdownConverter {
	return &MarkdownConverter{
			md: goldmark.New(
				goldmark.WithExtensions(
					extension.GFM, // 表格 / 任务列表等 GFM 语法
					mathjax.MathJax,
						highlighting.NewHighlighting(
						highlighting.WithStyle("github"), // 选一个内置主题名，后面可以改
						highlighting.WithFormatOptions(
							chromahtml.WithLineNumbers(true), // ✅ 行号
							chromahtml.WithClasses(true),     // 用 class，而不是 inline style，前端好控制
						),
					), // 数学公式支持
				),
				goldmark.WithParserOptions(
					parser.WithAutoHeadingID(),
				),
				goldmark.WithRendererOptions(
					gmhtml.WithHardWraps(),
					gmhtml.WithXHTML(),
				),
			),
		}
}

func (mc *MarkdownConverter) MdToHtml(markdown []byte) (string, error) {
	var buf bytes.Buffer
	if err := mc.md.Convert(markdown, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
