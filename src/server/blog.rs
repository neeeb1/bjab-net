use crate::blog::render::{render_post_html_from_md};
use crate::Post;
use axum::extract::{State, Path};
use axum::http::StatusCode;
use axum::response::{Html, IntoResponse};
use crate::AppState;
use askama::Template;

// Handlers for the route "/blog"
// Because of the blog mod, might want to rename - "/posts"?

#[derive(Template)]
#[template(path = "blog.html", escape = "none")]
struct BlogTemplate<'a> {
    posts: Vec<&'a Post>,
}

// Full list of blog posts
pub async fn get_blog_posts(State(state): State<AppState>) -> impl IntoResponse {
    let mut posts = state.posts_library.values().clone().collect::<Vec<&Post>>();
    posts.sort_by(|a, b| b.front_matter.date.cmp(&a.front_matter.date));

    let blog_template = BlogTemplate {
        posts
    };

    Html(
        blog_template
            .render()
            .expect("Failed to render from template"),
    )
}

// Retrieves post markdown post by slug and renders HTML page
// Routes for /blog/{post-slug}
pub async fn get_post(
    State(state): State<AppState>,
    Path(slug): Path<String>,
) -> impl IntoResponse {
    let response = match state.posts_library.get(&slug) {
        Some(post) => render_post_html_from_md(post.body.clone()),
        None => StatusCode::NOT_FOUND.to_string(),
    };

    Html(response)
}
