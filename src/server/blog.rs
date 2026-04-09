use crate::POSTS;
use crate::blog::render::render_html_from_md;
use axum::extract::Path;
use axum::http::StatusCode;
use axum::response::{Html, IntoResponse};

// Handler for the root of the site at "/blog"
// Because of the blog mod, might want to rename - "/posts"?

pub async fn get_post(Path(slug): Path<String>) -> impl IntoResponse {
    let response = match POSTS.get(&slug) {
        Some(post) => render_html_from_md(post.body.clone()),
        None => StatusCode::NOT_FOUND.to_string(),
    };

    Html(response)
}
