use std::collections::HashMap;

use std::sync::Arc;

use crate::blog::{Post, build_posts};

mod blog;
mod server;

// Axum App State, holds a HashMap of posts with post-slug as the key
// Allows for quick lookup during routing, and builds the library of Posts
// once at compile time
#[derive(Clone)]
struct AppState {
    posts_library: Arc<HashMap<String, Post>>,
}

#[tokio::main]
async fn main() {
    let posts_library: HashMap<String, Post> = build_posts()
        .expect("Failed to build vec of posts")
        .into_iter()
        .map(|post| (post.front_matter.slug.clone(), post))
        .collect();

    let state = AppState {
        posts_library: Arc::new(posts_library),
    };

    server::start_server(state).await
}
