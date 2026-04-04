# html

Renders HTML pages for the web server. The main entry point is `Html(db, path, writer)`, which dispatches to one of three page renderers based on the URL path:

- `ChannelsHtml` — lists all channels, sorted alphabetically, with a link and topic count for each
- `TopicsHtml` — lists all topics in a channel, sorted alphabetically, with a link and message count for each
- `MessagesHtml` — displays all messages in a channel/topic thread, sorted by message ID

HTML is written directly to an `io.StringWriter` without any template engine. User-supplied strings are escaped with Go's standard `html` package.
