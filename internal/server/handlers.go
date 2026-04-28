package server

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/neeeb1/bjab-net/internal/blog"
	"github.com/neeeb1/bjab-net/internal/projects"
)

var MAX_INDEX_POSTS = 3

var indexTemplate = template.Must(template.ParseFiles("web/templates/base.html", "web/templates/index.html"))
var blogTemplate = template.Must(template.ParseFiles("web/templates/base.html", "web/templates/blog.html"))
var postTemplate = template.Must(template.ParseFiles("web/templates/base.html", "web/templates/post.html"))
var projectsTemplate = template.Must(template.ParseFiles("web/templates/base.html", "web/templates/projects.html"))

type IndexData struct {
	Posts    []blog.Post
	Projects []projects.Project
}

type BlogData struct {
	Posts []blog.Post
}

type ProjectsData struct {
	Projects []projects.Project
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string
	}

	respBody := errorResponse{
		Error: msg,
	}

	data, _ := json.Marshal(respBody)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func (s *AppState) handleIndex(w http.ResponseWriter, r *http.Request) {
	posts, err := blog.SortedPosts(s.Posts)
	if err != nil {
		RespondWithError(w, 500, "Failed to sort blog posts")
		return
	}

	projects, err := projects.SortedProjects(s.Projects)
	if err != nil {
		RespondWithError(w, 500, "Failed to sort posts")
		return
	}

	if len(posts) > MAX_INDEX_POSTS {
		posts = posts[:MAX_INDEX_POSTS]
	}

	indexTemplate.ExecuteTemplate(w, "base", IndexData{Posts: posts, Projects: projects})
}

func (s *AppState) handlePost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	post, ok := s.Posts[slug]
	if !ok {
		RespondWithError(w, 404, "Post slug not found")
		return
	}
	postTemplate.ExecuteTemplate(w, "base", post.HTMLBody)
}

func (s *AppState) handleBlogList(w http.ResponseWriter, r *http.Request) {
	sortedPosts, err := blog.SortedPosts(s.Posts)
	if err != nil {
		RespondWithError(w, 500, "Failed to sort posts")
		return
	}
	blogTemplate.ExecuteTemplate(w, "base", BlogData{Posts: sortedPosts})
}

func (s *AppState) handleProjectList(w http.ResponseWriter, r *http.Request) {
	sortedProjects, err := projects.SortedProjects(s.Projects)
	if err != nil {
		RespondWithError(w, 500, "Failed to sort posts")
		return
	}
	projectsTemplate.ExecuteTemplate(w, "base", ProjectsData{Projects: sortedProjects})
}
