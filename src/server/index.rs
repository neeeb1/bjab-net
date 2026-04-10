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
    let index_template = IndexTemplate {
        posts: state.posts_library.values().clone().collect::<Vec<&Post>>(),
    };

    Html(
        index_template
            .render()
            .expect("Failed to render from template"),
    )
}
