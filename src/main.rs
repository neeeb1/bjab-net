use std::collections::HashMap;

use lazy_static::lazy_static;

use crate::blog::{Post, build_posts};

mod blog;
mod server;

// Lazy static, holds a HashMap of posts with post-slug as the key
// Allows for quick lookup during routing, and builds the library of Posts
// once at compile time
lazy_static! {
    static ref POSTS: HashMap<String, Post> = build_posts()
        .expect("Failed to build vec of posts")
        .into_iter()
        .map(|post| (post.front_matter.slug.clone(), post))
        .collect();
}

#[tokio::main]
async fn main() {
    server::start_server().await
}
