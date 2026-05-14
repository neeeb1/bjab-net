package server

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/neeeb1/bjab-net/internal/blog"
	"github.com/neeeb1/bjab-net/internal/meta"
	"github.com/neeeb1/bjab-net/internal/projects"
)

func TestMain(m *testing.M) {
	if err := os.Chdir("../.."); err != nil {
		panic("chdir to repo root: " + err.Error())
	}
	parseTemplates()
	os.Exit(m.Run())
}

func fakeState() AppState {
	posts := map[string]blog.Post{
		"first": {
			Metadata: meta.Metadata{Title: "First", Slug: "first", Date: "2024-01-01", Description: "first desc"},
			HTMLBody: template.HTML("<p>first body</p>"),
		},
		"second": {
			Metadata: meta.Metadata{Title: "Second", Slug: "second", Date: "2026-01-01", Description: "second desc"},
			HTMLBody: template.HTML("<p>second body</p>"),
		},
		"third": {
			Metadata: meta.Metadata{Title: "Third", Slug: "third", Date: "2025-01-01", Description: "third desc"},
			HTMLBody: template.HTML("<p>third body</p>"),
		},
		"fourth": {
			Metadata: meta.Metadata{Title: "Fourth", Slug: "fourth", Date: "2023-01-01", Description: "fourth desc"},
			HTMLBody: template.HTML("<p>fourth body</p>"),
		},
	}
	projs := map[string]projects.Project{
		"wasm-game": {
			Metadata:    meta.Metadata{Title: "WASM Game", Slug: "wasm-game", Date: "2025-06-01", Description: "g"},
			Embed:       "wasm",
			EmbedWidth:  1280,
			EmbedHeight: 720,
			HTMLBody:    template.HTML("<p>game</p>"),
		},
		"plain-proj": {
			Metadata: meta.Metadata{Title: "Plain", Slug: "plain-proj", Date: "2024-06-01"},
			HTMLBody: template.HTML("<p>plain</p>"),
		},
	}
	return AppState{
		Posts:        posts,
		Projects:     projs,
		WasmProjects: map[string]bool{"wasm-game": true},
	}
}

func doRequest(s AppState, method, path string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)
	req := httptest.NewRequest(method, path, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func TestHandleIndex_OKAndLimitsToMax(t *testing.T) {
	s := fakeState()
	rr := doRequest(s, "GET", "/")
	if rr.Code != 200 {
		t.Fatalf("status: %d", rr.Code)
	}
	body := rr.Body.String()

	// MAX_INDEX_POSTS = 3 — should include top 3 by date desc (second, third, first), not fourth.
	if !strings.Contains(body, "Second") {
		t.Error("missing Second post")
	}
	if !strings.Contains(body, "Third") {
		t.Error("missing Third post")
	}
	if !strings.Contains(body, "First") {
		t.Error("missing First post")
	}
	if strings.Contains(body, "Fourth") {
		t.Error("Fourth post should be excluded (over MAX_INDEX_POSTS)")
	}
	if !strings.Contains(body, "WASM Game") {
		t.Error("missing wasm project")
	}
}

func TestHandleBlogList_OK(t *testing.T) {
	s := fakeState()
	rr := doRequest(s, "GET", "/blog")
	if rr.Code != 200 {
		t.Fatalf("status: %d", rr.Code)
	}
	body := rr.Body.String()
	for _, want := range []string{"First", "Second", "Third", "Fourth"} {
		if !strings.Contains(body, want) {
			t.Errorf("blog list missing %q", want)
		}
	}
	// desc order — "Second" (2026) should appear before "Fourth" (2023)
	if strings.Index(body, "Second") > strings.Index(body, "Fourth") {
		t.Error("blog list not sorted desc")
	}
}

func TestHandlePost_KnownSlug(t *testing.T) {
	s := fakeState()
	rr := doRequest(s, "GET", "/blog/first")
	if rr.Code != 200 {
		t.Fatalf("status: %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "first body") {
		t.Errorf("post body missing in response")
	}
}

func TestHandlePost_UnknownSlugReturns404(t *testing.T) {
	s := fakeState()
	rr := doRequest(s, "GET", "/blog/does-not-exist")
	if rr.Code != 404 {
		t.Fatalf("status: got %d want 404", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("content-type: got %q", ct)
	}
	var body struct{ Error string }
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body.Error == "" {
		t.Error("expected non-empty Error field")
	}
}

func TestHandleProjectList_OK(t *testing.T) {
	s := fakeState()
	rr := doRequest(s, "GET", "/projects")
	if rr.Code != 200 {
		t.Fatalf("status: %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "WASM Game") || !strings.Contains(body, "Plain") {
		t.Errorf("projects list missing entries")
	}
}

func TestHandleProject_WasmSetsCOOPCOEP(t *testing.T) {
	s := fakeState()
	rr := doRequest(s, "GET", "/projects/wasm-game")
	if rr.Code != 200 {
		t.Fatalf("status: %d", rr.Code)
	}
	if got := rr.Header().Get("cross-origin-opener-policy"); got != "same-origin" {
		t.Errorf("COOP: got %q", got)
	}
	if got := rr.Header().Get("cross-origin-embedder-policy"); got != "require-corp" {
		t.Errorf("COEP: got %q", got)
	}
}

func TestHandleProject_PlainOmitsCOOPCOEP(t *testing.T) {
	s := fakeState()
	rr := doRequest(s, "GET", "/projects/plain-proj")
	if rr.Code != 200 {
		t.Fatalf("status: %d", rr.Code)
	}
	if rr.Header().Get("cross-origin-opener-policy") != "" {
		t.Error("COOP should not be set for non-wasm project")
	}
	if rr.Header().Get("cross-origin-embedder-policy") != "" {
		t.Error("COEP should not be set for non-wasm project")
	}
}

func TestHandleProject_UnknownSlugReturns404(t *testing.T) {
	s := fakeState()
	rr := doRequest(s, "GET", "/projects/nope")
	if rr.Code != 404 {
		t.Fatalf("status: got %d want 404", rr.Code)
	}
}

func TestHandleRSSFeed_ValidXMLAndCounts(t *testing.T) {
	s := fakeState()
	rr := doRequest(s, "GET", "/feed.xml")
	if rr.Code != 200 {
		t.Fatalf("status: %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/rss+xml" {
		t.Errorf("content-type: got %q", ct)
	}

	type rssItem struct {
		Title   string `xml:"title"`
		Link    string `xml:"link"`
		PubDate string `xml:"pubDate"`
	}
	type rssDoc struct {
		Items []rssItem `xml:"channel>item"`
	}
	var doc rssDoc
	if err := xml.Unmarshal(rr.Body.Bytes(), &doc); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}
	wantCount := len(s.Posts) + len(s.Projects)
	if len(doc.Items) != wantCount {
		t.Errorf("item count: got %d want %d", len(doc.Items), wantCount)
	}
}

func TestRespondWithError_ShapesJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	RespondWithError(rr, 418, "teapot")
	if rr.Code != 418 {
		t.Errorf("status: got %d", rr.Code)
	}
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("content-type: got %q", rr.Header().Get("Content-Type"))
	}
	var body struct{ Error string }
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body.Error != "teapot" {
		t.Errorf("Error: got %q", body.Error)
	}
}
