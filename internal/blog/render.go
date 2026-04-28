package blog

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var md = goldmark.New(
	goldmark.WithRendererOptions(
		renderer.WithNodeRenderers(util.Prioritized(NewBlogRenderer(), 100)),
	),
)

func RenderPost(body string) (string, error) {
	var buf bytes.Buffer
	err := md.Convert([]byte(body), &buf)
	return buf.String(), err
}
