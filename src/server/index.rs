use askama::Template;
use axum::response::{Html, IntoResponse};
use axum::extract::State;

use crate::{ AppState, blog::Post};

#[derive(Template)]
#[template(path = "index.html", escape = "none")]
struct IndexTemplate<'a> {
    posts: Vec<&'a Post>,
}

pub async fn build_index(State(state): State<AppState>) -> impl IntoResponse {
    let mut posts = state.posts_library.values().clone().collect::<Vec<&Post>>();
    posts.sort_by(|a, b| b.front_matter.date.cmp(&a.front_matter.date));

    let index_template = IndexTemplate {
        posts
    };

    Html(
        index_template
            .render()
            .expect("Failed to render from template"),
    )
}
