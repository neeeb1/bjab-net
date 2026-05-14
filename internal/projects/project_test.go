package projects

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleWasm = `---
title: WASM Game
date: 2025-08-01
slug: wasm-game
tags:
  - godot
description: A game.
embed: wasm
embed_width: 1280
embed_height: 720
---
Game body.
`

const samplePlain = `---
title: Plain
date: 2025-01-01
slug: plain
description: Plain project.
---
Body.
`

func writeProject(t *testing.T, root, slug, body string) {
	t.Helper()
	dir := filepath.Join(root, "web", "projects", slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "metadata.md"), []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestParseMetadata_EmbedFieldsPopulated(t *testing.T) {
	root := t.TempDir()
	writeProject(t, root, "wasm-game", sampleWasm)

	p, err := parseMetadata(filepath.Join(root, "web", "projects", "wasm-game", "metadata.md"))
	if err != nil {
		t.Fatalf("parseMetadata: %v", err)
	}
	if p.Metadata.Title != "WASM Game" {
		t.Errorf("title: got %q", p.Metadata.Title)
	}
	if p.Metadata.Slug != "wasm-game" {
		t.Errorf("slug: got %q", p.Metadata.Slug)
	}
	if p.Embed != "wasm" {
		t.Errorf("embed: got %q", p.Embed)
	}
	if p.EmbedWidth != 1280 || p.EmbedHeight != 720 {
		t.Errorf("embed dims: got %dx%d", p.EmbedWidth, p.EmbedHeight)
	}
	if !strings.Contains(string(p.HTMLBody), "Game body") {
		t.Errorf("HTMLBody missing: %q", p.HTMLBody)
	}
}

func TestParseMetadata_PlainProjectNoEmbed(t *testing.T) {
	root := t.TempDir()
	writeProject(t, root, "plain", samplePlain)
	p, err := parseMetadata(filepath.Join(root, "web", "projects", "plain", "metadata.md"))
	if err != nil {
		t.Fatalf("parseMetadata: %v", err)
	}
	if p.Embed != "" {
		t.Errorf("expected empty embed, got %q", p.Embed)
	}
}

func TestProject_GetDate(t *testing.T) {
	p := Project{}
	p.Metadata.Date = "2025-03-03"
	if p.GetDate() != "2025-03-03" {
		t.Errorf("GetDate: got %q", p.GetDate())
	}
}

func TestBuildProjects_KeysBySlug(t *testing.T) {
	root := t.TempDir()
	writeProject(t, root, "wasm-game", sampleWasm)
	writeProject(t, root, "plain", samplePlain)

	t.Chdir(root)

	got, err := BuildProjects()
	if err != nil {
		t.Fatalf("BuildProjects: %v", err)
	}
	if _, ok := got["wasm-game"]; !ok {
		t.Error("missing wasm-game")
	}
	if _, ok := got["plain"]; !ok {
		t.Error("missing plain")
	}
	if len(got) != 2 {
		t.Errorf("want 2, got %d", len(got))
	}
}

func TestSortedProjects_OrdersDesc(t *testing.T) {
	mk := func(slug, date string) Project {
		var p Project
		p.Metadata.Slug = slug
		p.Metadata.Date = date
		return p
	}
	in := map[string]Project{
		"a": mk("a", "2024-01-01"),
		"b": mk("b", "2026-01-01"),
		"c": mk("c", "2025-01-01"),
	}
	sorted, err := SortedProjects(in)
	if err != nil {
		t.Fatalf("SortedProjects: %v", err)
	}
	want := []string{"b", "c", "a"}
	for i, w := range want {
		if sorted[i].Metadata.Slug != w {
			t.Errorf("idx %d: got %q want %q", i, sorted[i].Metadata.Slug, w)
		}
	}
}
