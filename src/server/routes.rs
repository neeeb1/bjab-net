use crate::{
    AppState,
    server::{blog, index},
};
use axum::{Router, routing::get};
use tower_http::services::ServeDir;

// Returns router built with specified routes
// TODO: Add routes (blog/foo.html, projects/bar.html)

pub fn new_router(state: AppState) -> Router {
    Router::new()
        .route("/", get(index::build_index))
        .nest_service("/static", ServeDir::new("web/static"))
        .nest_service("/images", ServeDir::new("web/images"))
        .route("/blog",get(blog::get_blog_posts))
        .route("/blog/{slug}", get(blog::get_post))
        .with_state(state)
}
