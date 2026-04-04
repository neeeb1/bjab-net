use crate::server::routes::new_router;

pub mod blog;
pub mod routes;

pub async fn start_server() {
    let app = new_router();

    let listener = tokio::net::TcpListener::bind("0.0.0.0:1234").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
