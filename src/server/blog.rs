use crate::blog::{self, render::{TEST_BODY, render_html_from_md}};
use axum::response::{Html, IntoResponse};

// Handler for the root of the site at "/blog"
// Because of the blog mod, might want to rename - "/posts"?

pub async fn get_post() -> impl IntoResponse {
    let markdown_text = blog::read::read_markdown_file();
    Html(render_html_from_md(markdown_text)).into_response()
}
