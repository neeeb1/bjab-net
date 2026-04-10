use crate::{AppState, server::routes::new_router};

pub mod blog;
pub mod index;
pub mod routes;

pub async fn start_server(state: AppState) {
    let app = new_router(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:1234").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
