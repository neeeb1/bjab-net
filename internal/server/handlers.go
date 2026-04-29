package server

import (
	"encoding/json"
	"net/http"
	"sort"
	"text/template"
	"time"

	"github.com/neeeb1/bjab-net/internal/blog"
	"github.com/neeeb1/bjab-net/internal/projects"
)

var MAX_INDEX_POSTS = 3

var indexTemplate = template.Must(template.ParseFiles("web/templates/base.html", "web/templates/index.html"))
var blogTemplate = template.Must(template.ParseFiles("web/templates/base.html", "web/templates/blog.html"))
var postTemplate = template.Must(template.ParseFiles("web/templates/base.html", "web/templates/post.html"))
var projectsTemplate = template.Must(template.ParseFiles("web/templates/base.html", "web/templates/projects.html"))
var projectTemplate = template.Must(template.ParseFiles("web/templates/base.html", "web/templates/project.html"))
var feedTemplate = template.Must(template.ParseFiles("web/templates/feed.xml"))

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

type ProjectData struct {
	Project projects.Project
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

func (s *AppState) handleProject(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	project, ok := s.Projects[slug]
	if !ok {
		RespondWithError(w, 404, "Project not found")
		return
	}
	if project.Embed == "wasm" {
		w.Header().Set("cross-origin-opener-policy", "same-origin")
		w.Header().Set("cross-origin-embedder-policy", "require-corp")
	}
	projectTemplate.ExecuteTemplate(w, "base", ProjectData{Project: project})
}

func (s *AppState) handleProjectList(w http.ResponseWriter, r *http.Request) {
	sortedProjects, err := projects.SortedProjects(s.Projects)
	if err != nil {
		RespondWithError(w, 500, "Failed to sort posts")
		return
	}
	projectsTemplate.ExecuteTemplate(w, "base", ProjectsData{Projects: sortedProjects})
}

func (s *AppState) handleRSSFeed(w http.ResponseWriter, r *http.Request) {
	type FeedItem struct {
		Title       string
		Link        string
		PubDate     string
		Description string
	}

	type FeedData struct {
		Items []FeedItem
	}

	toRFC1123 := func(date string) string {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			return ""
		}
		return t.Format(time.RFC1123Z)
	}

	items := make([]FeedItem, 0, len(s.Posts)+len(s.Projects))
	for _, p := range s.Posts {
		items = append(items, FeedItem{
			Title:       p.Metadata.Title,
			Link:        "https://bjab.net/blog/" + p.Metadata.Slug,
			PubDate:     toRFC1123(p.Metadata.Date),
			Description: p.Metadata.Description,
		})
	}
	for _, p := range s.Projects {
		items = append(items, FeedItem{
			Title:       p.Metadata.Title,
			Link:        "https://bjab.net/projects/" + p.Metadata.Slug,
			PubDate:     toRFC1123(p.Metadata.Date),
			Description: p.Metadata.Description,
		})
	}

	var err error
	sort.SliceStable(items, func(i, j int) bool {
		t1, e1 := time.Parse(time.RFC1123Z, items[i].PubDate)
		t2, e2 := time.Parse(time.RFC1123Z, items[j].PubDate)
		if e1 != nil {
			err = e1
		}
		if e2 != nil {
			err = e2
		}
		return t1.After(t2)
	})
	if err != nil {
		RespondWithError(w, 500, "Failed to sort articles")
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	feedTemplate.Execute(w, FeedData{Items: items})
}
