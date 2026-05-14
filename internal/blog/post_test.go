package blog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/neeeb1/bjab-net/internal/meta"
)

func metaFor(slug, date string) meta.Metadata {
	return meta.Metadata{Slug: slug, Date: date}
}

const sampleValid = `---
title: Hello World
date: 2025-06-15
slug: hello-world
tags:
  - go
  - web
description: A first post.
draft: false
---
# Heading

Some body text.
`

const sampleDraft = `---
title: Draft Post
date: 2025-07-01
slug: draft-post
draft: true
---
Body.
`

const sampleMinimal = `---
title: Bare
date: 2025-01-01
slug: bare
---
Body only.
`

func writePost(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return p
}

func TestParseMarkdownFile_ValidFrontmatter(t *testing.T) {
	dir := t.TempDir()
	path := writePost(t, dir, "p.md", sampleValid)

	got, err := parseMarkdownFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Metadata.Title != "Hello World" {
		t.Errorf("title: got %q", got.Metadata.Title)
	}
	if got.Metadata.Slug != "hello-world" {
		t.Errorf("slug: got %q", got.Metadata.Slug)
	}
	if got.Metadata.Date != "2025-06-15" {
		t.Errorf("date: got %q", got.Metadata.Date)
	}
	if len(got.Metadata.Tags) != 2 || got.Metadata.Tags[0] != "go" {
		t.Errorf("tags: got %+v", got.Metadata.Tags)
	}
	if got.Metadata.Draft {
		t.Error("draft: expected false")
	}
	if !strings.Contains(string(got.HTMLBody), "Heading") {
		t.Errorf("HTMLBody missing heading content: %q", got.HTMLBody)
	}
	if !strings.Contains(got.MdBody, "# Heading") {
		t.Errorf("MdBody missing raw markdown: %q", got.MdBody)
	}
}

func TestParseMarkdownFile_DraftFlagParsed(t *testing.T) {
	dir := t.TempDir()
	path := writePost(t, dir, "d.md", sampleDraft)
	got, err := parseMarkdownFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Metadata.Draft {
		t.Error("expected Draft=true")
	}
}

func TestParseMarkdownFile_MissingFile(t *testing.T) {
	_, err := parseMarkdownFile(filepath.Join(t.TempDir(), "nope.md"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestPost_GetDate(t *testing.T) {
	p := Post{}
	p.Metadata.Date = "2025-01-01"
	if p.GetDate() != "2025-01-01" {
		t.Errorf("GetDate: got %q", p.GetDate())
	}
}

func TestBuildPosts_SkipsDraftsAndKeysBySlug(t *testing.T) {
	root := t.TempDir()
	postDir := filepath.Join(root, "web", "posts")
	if err := os.MkdirAll(postDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writePost(t, postDir, "valid.md", sampleValid)
	writePost(t, postDir, "draft.md", sampleDraft)
	writePost(t, postDir, "min.md", sampleMinimal)

	t.Chdir(root)

	posts, err := BuildPosts()
	if err != nil {
		t.Fatalf("BuildPosts: %v", err)
	}
	if _, ok := posts["draft-post"]; ok {
		t.Error("expected draft to be skipped")
	}
	if _, ok := posts["hello-world"]; !ok {
		t.Error("expected hello-world present")
	}
	if _, ok := posts["bare"]; !ok {
		t.Error("expected bare present")
	}
	if len(posts) != 2 {
		t.Errorf("want 2 posts, got %d", len(posts))
	}
}

func TestSortedPosts_OrdersDesc(t *testing.T) {
	posts := map[string]Post{
		"a": {Metadata: metaFor("a", "2024-01-01")},
		"b": {Metadata: metaFor("b", "2026-01-01")},
		"c": {Metadata: metaFor("c", "2025-01-01")},
	}
	sorted, err := SortedPosts(posts)
	if err != nil {
		t.Fatalf("SortedPosts: %v", err)
	}
	want := []string{"b", "c", "a"}
	for i, w := range want {
		if sorted[i].Metadata.Slug != w {
			t.Errorf("idx %d: got %q want %q", i, sorted[i].Metadata.Slug, w)
		}
	}
}
