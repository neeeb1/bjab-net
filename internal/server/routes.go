package server

import (
	"net/http"

	"github.com/neeeb1/bjab-net/internal/projects"
)

func RegisterRoutes(mux *http.ServeMux, state AppState) {
	mux.HandleFunc("GET /", state.handleIndex)

	mux.HandleFunc("GET /blog", state.handleBlogList)
	mux.HandleFunc("GET /blog/{slug}", state.handlePost)

	mux.HandleFunc("GET /projects", state.handleProjectList)
	mux.HandleFunc("GET /projects/{slug}", state.handleProject)

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /images/", http.StripPrefix("/images/", http.FileServer(http.Dir("web/images"))))

	projectsHandler := http.StripPrefix("/projects/", http.FileServer(http.Dir("web/projects")))
	mux.Handle("GET /projects/", projects.WasmHeaders(state.WasmProjects, projectsHandler))
}
