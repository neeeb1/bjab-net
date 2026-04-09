use crate::server::blog;
use axum::{Router, routing::get};
use tower_http::services::{ServeDir, ServeFile};

// Returns router built with specified routes
// TODO: Add routes (blog/foo.html, projects/bar.html)

pub fn new_router() -> Router {
    Router::new()
        .route_service("/", ServeFile::new("web/index.html"))
        .nest_service("/static", ServeDir::new("web/static"))
        .route("/blog/{slug}", get(blog::get_post))
}
