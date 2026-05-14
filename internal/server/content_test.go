package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/neeeb1/bjab-net/internal/blog"
	"github.com/neeeb1/bjab-net/internal/meta"
	"github.com/neeeb1/bjab-net/internal/projects"
	"go.yaml.in/yaml/v2"
)

type postFrontmatter struct {
	meta.Metadata `yaml:",inline"`
}

type projectFrontmatter struct {
	meta.Metadata `yaml:",inline"`
	Embed         string `yaml:"embed"`
	EmbedWidth    int    `yaml:"embed_width"`
	EmbedHeight   int    `yaml:"embed_height"`
}

func parseFrontmatter(t *testing.T, label, path string, dst any) bool {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("%s: read: %v", label, err)
		return false
	}
	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) < 3 {
		t.Errorf("%s: missing frontmatter delimiters (expected `---` ... `---`)", label)
		return false
	}
	if err := yaml.Unmarshal([]byte(parts[1]), dst); err != nil {
		t.Errorf("%s: yaml parse: %v", label, err)
		return false
	}
	return true
}

func validateMeta(t *testing.T, label string, m meta.Metadata) {
	t.Helper()
	if strings.TrimSpace(m.Title) == "" {
		t.Errorf("%s: missing title", label)
	}
	if strings.TrimSpace(m.Slug) == "" {
		t.Errorf("%s: missing slug", label)
	}
	if strings.TrimSpace(m.Date) == "" {
		t.Errorf("%s: missing date", label)
	} else if _, err := time.Parse("2006-01-02", m.Date); err != nil {
		t.Errorf("%s: date %q does not match YYYY-MM-DD: %v", label, m.Date, err)
	}
}

func TestRealContent_PostsHaveValidFrontmatter(t *testing.T) {
	entries, err := os.ReadDir("web/posts")
	if err != nil {
		t.Fatalf("read web/posts: %v", err)
	}
	seenSlug := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		var fm postFrontmatter
		path := filepath.Join("web/posts", e.Name())
		if !parseFrontmatter(t, e.Name(), path, &fm) {
			continue
		}
		validateMeta(t, e.Name(), fm.Metadata)
		if fm.Slug != "" {
			if prev, dup := seenSlug[fm.Slug]; dup {
				t.Errorf("duplicate slug %q in %s (also in %s) — map would silently overwrite", fm.Slug, e.Name(), prev)
			}
			seenSlug[fm.Slug] = e.Name()
		}
	}
}

func TestRealContent_ProjectsHaveValidFrontmatter(t *testing.T) {
	entries, err := os.ReadDir("web/projects")
	if err != nil {
		t.Fatalf("read web/projects: %v", err)
	}
	seenSlug := map[string]string{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		metaPath := filepath.Join("web/projects", e.Name(), "metadata.md")
		if _, err := os.Stat(metaPath); err != nil {
			t.Errorf("%s: missing metadata.md", e.Name())
			continue
		}
		var fm projectFrontmatter
		if !parseFrontmatter(t, e.Name(), metaPath, &fm) {
			continue
		}
		validateMeta(t, e.Name(), fm.Metadata)
		if fm.Slug != "" {
			if prev, dup := seenSlug[fm.Slug]; dup {
				t.Errorf("duplicate slug %q in %s (also in %s)", fm.Slug, e.Name(), prev)
			}
			seenSlug[fm.Slug] = e.Name()
		}

		if fm.Embed == "wasm" {
			files, err := os.ReadDir(filepath.Join("web/projects", e.Name()))
			if err != nil {
				t.Errorf("%s: read dir: %v", e.Name(), err)
				continue
			}
			hasWasm := false
			for _, f := range files {
				if strings.HasSuffix(f.Name(), ".wasm") {
					hasWasm = true
					break
				}
			}
			if !hasWasm {
				t.Errorf("%s: embed=wasm but no .wasm file in project dir", e.Name())
			}
		}
	}
}

func loadRealState(t *testing.T) AppState {
	t.Helper()
	posts, err := blog.BuildPosts()
	if err != nil {
		t.Fatalf("BuildPosts: %v", err)
	}
	projs, err := projects.BuildProjects()
	if err != nil {
		t.Fatalf("BuildProjects: %v", err)
	}
	wasm, err := projects.FindWasmProjects()
	if err != nil {
		t.Fatalf("FindWasmProjects: %v", err)
	}
	return AppState{Posts: posts, Projects: projs, WasmProjects: wasm}
}

func TestRealContent_BuildsWithoutError(t *testing.T) {
	_ = loadRealState(t)
}

func TestRealContent_AllRoutesRespond200(t *testing.T) {
	s := loadRealState(t)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	check := func(path string) {
		t.Helper()
		req := httptest.NewRequest("GET", path, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != 200 {
			t.Errorf("GET %s: status %d, body: %s", path, rr.Code, rr.Body.String())
		}
	}

	check("/")
	check("/blog")
	check("/projects")
	check("/feed.xml")
	for slug := range s.Posts {
		check("/blog/" + slug)
	}
	for slug := range s.Projects {
		check("/projects/" + slug)
	}
}
