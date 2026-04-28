package projects

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func FindWasmProjects() (map[string]bool, error) {
	projects := map[string]bool{}
	entries, err := os.ReadDir("web/projects")
	if err != nil {
		return projects, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		files, err := os.ReadDir(filepath.Join("web/projects", entry.Name()))
		if err != nil {
			return projects, err
		}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".wasm") {
				projects[entry.Name()] = true
				break
			}
		}
	}

	return projects, nil
}

func WasmHeaders(WasmProjects map[string]bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/projects/"), "/", 2)
		if len(parts) > 0 && WasmProjects[parts[0]] {
			w.Header().Set("cross-origin-opener-policy", "same-origin")
			w.Header().Set("cross-origin-embedder-policy", "require-corp")
		}
		next.ServeHTTP(w, r)
	})
}
