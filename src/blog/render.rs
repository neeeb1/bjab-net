use askama::Template;
use comrak::{Arena, Options, parse_document};

use super::format::BlogFormatter;
// TODO: [X] Generate HTML from templates and .md files
// TODO: [X] Add comrak formatter (https://docs.rs/comrak/latest/comrak/macro.create_formatter.html)
// TODO: [X] Create posts struct to hold post meta data + body
// TODO: [X] Wire up greymatter with comrak to get post meta data (slug, title, description, etc.)
// TODO: [ ] Write all html at startup rather than each time handler is called?

#[derive(Template)]
#[template(path = "post.html", escape = "none")]
struct PostTemplate<'a> {
    post_body: &'a str,
}

pub fn render_html_from_md(markdown_body: String) -> String {
    let options = Options::default();
    let arena = Arena::new();
    let doc = parse_document(&arena, &markdown_body, &options);

    let mut body = String::new();

    BlogFormatter::format_document(doc, &options, &mut body).unwrap();

    let post_template = PostTemplate { post_body: &body };

    post_template.render().expect("Failed to render from template")
}
