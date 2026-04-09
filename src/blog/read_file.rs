use super::{FrontMatter, Post};
use gray_matter::engine::YAML;
use gray_matter::{Matter, ParsedEntity, Result};
use std::{fs::File, io::Read};

pub fn read_markdown_file() -> Result<Post> {
    let mut post = Post::default();
    let mut file = File::open("posts/test_post.md").expect("Unable to open file");
    let mut contents = String::new();

    file.read_to_string(&mut contents)
        .expect("Unable to read file");

    let matter = Matter::<YAML>::new();
    let parsed: ParsedEntity = matter.parse(&contents)?;

    let front_matter: FrontMatter = parsed.data.as_ref().unwrap().deserialize()?;

    post.front_matter = front_matter;
    post.body = parsed.content;
    Ok(post)
}
