# server_types

Defines the Go structs that mirror the JSON types returned by the Zulip REST API:

- `ServerSubscription` — a channel the user is subscribed to (stream ID and name)
- `ServerMessage` — a single message (content, sender, stream, topic, etc.)

The field names use Go camelCase conventions with `json:` tags mapping to the snake_case names used in the Zulip API responses.
