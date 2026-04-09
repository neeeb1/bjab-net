use std::collections::HashMap;

use lazy_static::lazy_static;

use crate::blog::{Post, build_posts};

mod blog;
mod server;

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
