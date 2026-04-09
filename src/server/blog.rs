use crate::blog::{self, render::render_html_from_md};
use axum::response::{Html, IntoResponse};

// Handler for the root of the site at "/blog"
// Because of the blog mod, might want to rename - "/posts"?

pub async fn get_post() -> impl IntoResponse {
    let post =
        blog::read_file::read_markdown_file().expect("Failed to parse post struct from markdown");
    Html(render_html_from_md(post.body)).into_response()
}
