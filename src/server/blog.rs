use crate::blog::posts;
use axum::response::{Html, IntoResponse};

// Handler for the root of the site at "/blog"
// Because of the blog mod, might want to rename - "/posts"?

pub async fn get_post() -> impl IntoResponse {
    Html(posts::render()).into_response()
}
