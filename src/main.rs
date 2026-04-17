use std::collections::{HashMap, HashSet};

use std::sync::Arc;

use crate::blog::{Post, build_posts};
use crate::projects::wasm::find_wasm_projects;

mod blog;
mod projects;
mod server;

// Axum App State, holds a HashMap of posts with post-slug as the key
// Allows for quick lookup during routing, and builds the library of Posts
// once at compile time
#[derive(Clone)]
struct AppState {
    posts_library: Arc<HashMap<String, Post>>,
    wasm_projects: Arc<HashSet<String>>,
}

#[tokio::main]
async fn main() {
    let posts_library: HashMap<String, Post> = build_posts()
        .expect("Failed to build vec of posts")
        .into_iter()
        .map(|post| (post.front_matter.slug.clone(), post))
        .collect();

    let wasm_projects = find_wasm_projects();

    let state = AppState {
        posts_library: Arc::new(posts_library),
        wasm_projects: Arc::new(wasm_projects),
    };

    server::start_server(state).await
}
