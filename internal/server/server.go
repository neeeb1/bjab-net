package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neeeb1/bjab-net/internal/blog"
	"github.com/neeeb1/bjab-net/internal/projects"
)

var PORT = 3000

type AppState struct {
	Posts        map[string]blog.Post
	Projects     map[string]projects.Project
	WasmProjects map[string]bool
}

func StartServer() {
	mux := http.NewServeMux()
	state := AppState{}

	var err error
	state.Posts, err = blog.BuildPosts()
	if err != nil {
		log.Fatal("Failed to build list of blog posts: ", err)
	}

	state.Projects, err = projects.BuildProjects()
	if err != nil {
		log.Fatal("Failed to build list of projects: ", err)
	}

	state.WasmProjects, err = projects.FindWasmProjects()
	if err != nil {
		log.Fatal("Failed to build list of wasm projects")
	}

	RegisterRoutes(mux, state)

	log.Printf("Server starting, listening on port :%d", PORT)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", PORT), mux)
}
