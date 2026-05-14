package projects

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFindWasmProjects_DetectsWasmFile(t *testing.T) {
	root := t.TempDir()
	projDir := filepath.Join(root, "web", "projects")
	if err := os.MkdirAll(filepath.Join(projDir, "game-a"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(projDir, "post-only"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projDir, "game-a", "build.wasm"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projDir, "post-only", "metadata.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Chdir(root)

	got, err := FindWasmProjects()
	if err != nil {
		t.Fatalf("FindWasmProjects: %v", err)
	}
	if !got["game-a"] {
		t.Error("expected game-a to be detected")
	}
	if got["post-only"] {
		t.Error("did not expect post-only to be detected")
	}
}

func TestFindWasmProjects_EmptyDir(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "web", "projects"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	got, err := FindWasmProjects()
	if err != nil {
		t.Fatalf("FindWasmProjects: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("want empty map, got %+v", got)
	}
}

func TestWasmHeaders_SetsCOOPCOEPForWasmSlug(t *testing.T) {
	wasm := map[string]bool{"game-a": true}
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := WasmHeaders(wasm, inner)

	req := httptest.NewRequest("GET", "/projects/game-a/index.html", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if got := rr.Header().Get("cross-origin-opener-policy"); got != "same-origin" {
		t.Errorf("COOP: got %q", got)
	}
	if got := rr.Header().Get("cross-origin-embedder-policy"); got != "require-corp" {
		t.Errorf("COEP: got %q", got)
	}
}

func TestWasmHeaders_NoHeadersForNonWasmSlug(t *testing.T) {
	wasm := map[string]bool{"game-a": true}
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := WasmHeaders(wasm, inner)

	req := httptest.NewRequest("GET", "/projects/plain/index.html", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Header().Get("cross-origin-opener-policy") != "" {
		t.Error("COOP should not be set for non-wasm slug")
	}
	if rr.Header().Get("cross-origin-embedder-policy") != "" {
		t.Error("COEP should not be set for non-wasm slug")
	}
}

func TestWasmHeaders_CallsNext(t *testing.T) {
	called := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusTeapot)
	})
	h := WasmHeaders(map[string]bool{}, inner)
	req := httptest.NewRequest("GET", "/projects/anything", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if !called {
		t.Error("next handler not invoked")
	}
	if rr.Code != http.StatusTeapot {
		t.Errorf("status: got %d", rr.Code)
	}
}
