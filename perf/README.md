# perf

A standalone binary that serves the go-zulip web server and also contains performance testing harnesses.

## Web server

`webServer()` fetches live data from Zulip via the `zulip` package and starts an HTTP server on port 8080. CSS is served from `/style.css` (embedded in the binary at build time). Every HTML page is wrapped with a `<head>` that links the stylesheet.

## Performance testing

`buildBigDb()` generates a large synthetic database (20 channels × 20 topics × 1000 message batches) for use in benchmarking the HTML rendering code without needing a live Zulip connection. The `channels()` and `topicsAndMessages()` functions run rendering loops and report throughput.

To run the web server:
```bash
go run ./perf/
```
