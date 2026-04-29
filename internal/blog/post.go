package blog

import (
	"html/template"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/neeeb1/bjab-net/internal/meta"
	"go.yaml.in/yaml/v2"
)

type Post struct {
	Metadata meta.Metadata
	MdBody   string
	HTMLBody template.HTML
}

func (p Post) GetDate() string { return p.Metadata.Date }

func parseMarkdownFile(path string) (Post, error) {
	var result Post
	data, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}

	// split --- divider indicating yaml frontmater
	parts := strings.SplitN(string(data), "---", 3)

	// parts[1] = YAML frontmater, parts[2] = markdown body
	var m meta.Metadata
	yaml.Unmarshal([]byte(parts[1]), &m)

	// Render part[2] (body) to valid html
	markdown := strings.TrimSpace(parts[2])
	html, err := RenderPost(markdown)
	if err != nil {
		return result, err
	}

	return Post{Metadata: m, MdBody: markdown, HTMLBody: template.HTML(html)}, err
}

func BuildPosts() (map[string]Post, error) {
	posts := make(map[string]Post)

	entries, err := os.ReadDir("web/posts")
	if err != nil {
		return posts, err
	}

	for _, e := range entries {
		post, err := parseMarkdownFile(filepath.Join("web/posts", e.Name()))
		if err != nil {
			return posts, err
		}
		if post.Metadata.Draft == true {
			continue
		}
		posts[post.Metadata.Slug] = post
	}

	return posts, nil
}

func SortedPosts(posts map[string]Post) ([]Post, error) {
	return meta.SortByDate(slices.Collect(maps.Values(posts)))
}
