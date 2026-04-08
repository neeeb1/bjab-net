use askama::Template;
use comrak::{Arena, Options, markdown_to_html, parse_document};
use crate::blog::read;

use super::format::BlogFormatter;
// TODO: [X] Generate HTML from templates and .md files
// TODO: [X] Add comrak formatter (https://docs.rs/comrak/latest/comrak/macro.create_formatter.html)
// TODO: [ ] Create posts struct to hold post meta data + body
// TODO: [ ] Wire up greymatter with comrak to get post meta data (slug, title, description, etc.)
// TODO: [ ] Write all html at startup rather than each time handler is called?

#[derive(Template)]
#[template(path = "post.html", escape = "none")]
struct PostTemplate<'a> {
    post_body: &'a str,
}

pub const TEST_BODY: &str = "# Title of blog
## Subtitle of blog: and why it's important
Date and Time of publishing
Last edit: at a time

![image alt text](/images/title_image.png)

This is the first paragraph. You can tell, because it's under the title image and before the other sections.

### Section 2, another section

This is the second paragraph. You can tell, because it's under the second section (###) heading image and before the other sections. It's also after the first section.";

pub fn render_html_from_md(markdown_body: String) -> String {
    let options = Options::default();
    let arena = Arena::new();
    let doc = parse_document(&arena, &markdown_body, &options);

    let mut body = String::new();
    
    
    BlogFormatter::format_document(doc, &options, &mut body).unwrap();

    let post_template = PostTemplate { post_body: &body };

    post_template.render().unwrap()
}
