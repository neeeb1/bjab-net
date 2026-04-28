package projects

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/neeeb1/bjab-net/internal/blog"
	"go.yaml.in/yaml/v2"
)

type Project struct {
	Metadata blog.Metadata
}

func parseMetadata(path string) (Project, error) {
	var result Project
	data, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}

	// split --- divider indicating yaml frontmater
	parts := strings.SplitN(string(data), "---", 3)

	// parts[1] = YAML frontmater, parts[2] = markdown body
	var meta blog.Metadata
	yaml.Unmarshal([]byte(parts[1]), &meta)

	return Project{Metadata: meta}, err
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
	result := make([]Project, 0, len(projects))

	for _, p := range projects {
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
