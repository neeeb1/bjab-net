use std::{fs::File, io::Read};

pub fn read_markdown_file() -> String {
    let mut file = File::open("posts/test_post.md").expect("Unable to open file");
    let mut contents = String::new();

    file.read_to_string(&mut contents)
        .expect("Unable to read file");
    contents
}
