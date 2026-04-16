use std::{collections::HashSet, fs::read_dir};
use axum::{
    Router,
    extract::{Request, State},
    http::{HeaderName, HeaderValue},
    middleware::Next,
    response::Response,
    routing::get,
};

use crate::AppState;

pub fn find_wasm_projects() -> HashSet<String> {
    let mut wasm_projects = HashSet::new();
    let Ok(projects) = read_dir("web/projects") else {
        return wasm_projects;
    };

    for project in projects.flatten() {
        let path = project.path();
        if !path.is_dir() {
            continue;
        }

        let has_wasm = read_dir(&path)
            .into_iter()
            .flatten()
            .flatten()
            .any(|f| f.path().extension().is_some_and(|e| e == "wasm"));

        if has_wasm {
            if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
                wasm_projects.insert(name.to_string());
            }
        }
    }
    wasm_projects
}

pub async fn wasm_headers(State(state): State<AppState>, req: Request, next: Next) -> Response {
    let project_slug = req
        .uri()
        .path()
        .strip_prefix("/projects/")
        .and_then(|rest| rest.split('/').next());

    let needs_headers = project_slug.is_some_and(|slug| state.wasm_projects.contains(slug));

    let mut response = next.run(req).await;

    if needs_headers {
        let h = response.headers_mut();
        h.insert(
            HeaderName::from_static("cross-origin-opener-policy"),
            HeaderValue::from_static("same-origin"),
        );
        h.insert(
            HeaderName::from_static("cross-origin-embedder-policy"),
            HeaderValue::from_static("require-corp"),
        );
    }
    response
}