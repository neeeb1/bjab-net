package blog

import (
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.yaml.in/yaml/v2"
)

type Metadata struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Slug        string   `yaml:"slug"`
	Tags        []string `yaml:"tags"`
	Description string   `yaml:"description"`
	Draft       bool     `yaml:"draft"`
}

type Post struct {
	Metadata Metadata
	MdBody   string
	HTMLBody template.HTML
}

func parseMarkdownFile(path string) (Post, error) {
	var result Post
	data, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}

	// split --- divider indicating yaml frontmater
	parts := strings.SplitN(string(data), "---", 3)

	// parts[1] = YAML frontmater, parts[2] = markdown body
	var meta Metadata
	yaml.Unmarshal([]byte(parts[1]), &meta)

	// Render part[2] (body) to valid html
	markdown := strings.TrimSpace(parts[2])
	html, err := RenderPost(markdown)
	if err != nil {
		return result, err
	}

	return Post{Metadata: meta, MdBody: markdown, HTMLBody: template.HTML(html)}, err
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
		posts[post.Metadata.Slug] = post
	}

	return posts, nil
}

func SortedPosts(posts map[string]Post) ([]Post, error) {
	result := make([]Post, 0, len(posts))

	for _, p := range posts {
		result = append(result, p)
	}

	var err error
	sort.SliceStable(result, func(i int, j int) bool {
		layout := "2006-01-02"

		t1, e1 := time.Parse(layout, result[i].Metadata.Date)
		t2, e2 := time.Parse(layout, result[j].Metadata.Date)
		if e1 != nil {
			err = e1
		}
		if e2 != nil {
			err = e2
		}

		return t1.After(t2)
	})
	return result, err
}
