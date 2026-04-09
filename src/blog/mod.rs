use std::fs::read_dir;

use serde::Deserialize;

use crate::blog::read_file::read_markdown_file;

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
    pub title: String,
    pub date: String,
    pub slug: String,
    pub tags: Vec<String>,
    pub description: String,
}

pub fn build_posts() -> Result<Vec<Post>, std::io::Error> {
    let mut posts = Vec::<Post>::new();

    let markdown_files = read_dir("posts")?;
    for entry in markdown_files {
        let entry = entry?;
        let path = entry.path();
        posts.push(read_markdown_file(path).expect("Failed to parse post struct from markdown"));
    }

    Ok(posts)
}
