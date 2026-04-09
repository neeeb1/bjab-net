use askama::Template;
use axum::response::{Html, IntoResponse};

use crate::{POSTS, blog::Post};

#[derive(Template)]
#[template(path = "index.html", escape = "none")]
struct IndexTemplate<'a> {
    posts: Vec<&'a Post>,
}

pub async fn build_index() -> impl IntoResponse {
    let index_template = IndexTemplate{ posts: POSTS.values().clone().collect::<Vec<&Post>>()};
    
    Html(index_template.render().expect("Failed to render from template"))
}