package projects

import (
	"html/template"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/neeeb1/bjab-net/internal/blog"
	"github.com/neeeb1/bjab-net/internal/meta"
	"go.yaml.in/yaml/v2"
)

type projectMetadata struct {
	meta.Metadata `yaml:",inline"`
	Embed         string `yaml:"embed"`
}

type Project struct {
	Metadata meta.Metadata
	Embed    string
	MdBody   string
	HTMLBody template.HTML
}

func (p Project) GetDate() string { return p.Metadata.Date }

func parseMetadata(path string) (Project, error) {
	var result Project
	data, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}

	// split --- divider indicating yaml frontmater
	parts := strings.SplitN(string(data), "---", 3)

	// parts[1] = YAML frontmater, parts[2] = markdown body
	var m projectMetadata
	yaml.Unmarshal([]byte(parts[1]), &m)

	// read and render markdown body (parts[2])
	md := strings.TrimSpace(parts[2])
	html, err := blog.RenderPost(md)
	if err != nil {
		return result, err
	}

	return Project{Metadata: m.Metadata, Embed: m.Embed, MdBody: md, HTMLBody: template.HTML(html)}, err
}

func BuildProjects() (map[string]Project, error) {
	projects := make(map[string]Project)

	entries, err := os.ReadDir("web/projects")
	if err != nil {
		return projects, err
	}

	for _, e := range entries {
		project, err := parseMetadata(filepath.Join("web/projects/", e.Name(), "metadata.md"))
		if err != nil {
			return projects, err
		}
		projects[project.Metadata.Slug] = project
	}

	return projects, nil
}

func SortedProjects(projects map[string]Project) ([]Project, error) {
	return meta.SortByDate(slices.Collect(maps.Values(projects)))
}
