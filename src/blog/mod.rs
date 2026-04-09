use serde::Deserialize;

mod format;
pub mod read_file;
pub mod render;

#[derive(Default, Debug)]
pub struct Post {
    pub front_matter: FrontMatter,
    pub body: String,
}

#[derive(Default, Deserialize, Debug)]
pub struct FrontMatter {
    title: String,
    date: String,
    slug: String,
    tags: Vec<String>,
    description: String,
}