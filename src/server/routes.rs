use crate::{
    AppState,
    projects::wasm::{self, wasm_headers},
    server::{blog, index},
};
use axum::{
    Router,
    extract::{Request, State},
    http::{HeaderName, HeaderValue},
    middleware::{self, Next},
    response::Response,
    routing::get,
};
use tower_http::services::ServeDir;

// Returns router built with specified routes
// TODO: Add routes (blog/foo.html, projects/bar.html)

pub fn new_router(state: AppState) -> Router {
    let projects = Router::new()
        .fallback_service(ServeDir::new("web/projects"))
        .layer(middleware::from_fn_with_state(state.clone(), wasm_headers));

    Router::new()
        .route("/", get(index::build_index))
        .nest_service("/static", ServeDir::new("web/static"))
        .nest_service("/images", ServeDir::new("web/images"))
        .route("/blog", get(blog::get_blog_posts))
        .route("/blog/{slug}", get(blog::get_post))
        .nest("/projects", projects)
        .with_state(state)
}
