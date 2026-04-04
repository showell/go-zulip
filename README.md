# go-zulip

A Go web server that fetches and displays Zulip channel/topic/message data.

## Setup

Copy `config.json.example` to `config.json` and fill in your credentials:

```bash
cp config.json.example config.json
```

Edit `config.json`:

```json
{
    "email": "you@example.com",
    "api_key": "your_api_key_here",
    "base_url": "https://yourorg.zulipchat.com"
}
```

Your API key can be found in Zulip under **Settings → Account & Privacy → API key**.

`config.json` is gitignored and will never be committed.

## Running the web server

```bash
go run ./perf/
```

The server starts on port 8080. Visit `http://localhost:8080/channels`.
