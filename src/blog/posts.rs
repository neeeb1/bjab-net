use comrak::{Options, markdown_to_html};
// TODO: Generate HTML from templates and .md files

pub fn render() -> String {
    let body = markdown_to_html("# Title of blog
## Subtitle of blog: and why it's important
Date and Time of publishing
Last edit: at a time

![image alt text](/images/title_image.png)

This is the first paragraph. You can tell, because it's under the title image and before the other sections.

### Section 2, another section

This is the second paragraph. You can tell, because it's under the second section (###) heading image and before the other sections. It's also after the first section.", &Options::default());

    let header = "
    <html>
    <head>
    <title>Hello, web</title>
    </head>
    ";

    let footer = "
    </html>";

    header.to_owned() + &body + footer
}
