# zulip

Fetches data from the Zulip REST API and populates a `database.Database`.

`BuildDatabase(configPath)` is the main entry point. It reads credentials from a JSON config file, then:

1. Fetches all subscribed channels
2. Fetches the newest 1000 messages
3. Backfills older messages in batches of 5000, up to a limit of 50,000

Authentication uses HTTP Basic Auth with an email and API key. See the root `README.md` for config setup.

## Current limitations

The fetching is a one-time snapshot that happens at server startup. Once the database is built, it reflects the state of the Zulip server at that moment and does not update. New messages posted while the server is running will not appear until the server is restarted.

## Future plans

The goal is to transition from a static snapshot to a live, continuously-updated view. Zulip provides a real-time event API that allows a client to register for events and receive notifications as new messages are posted, channels are created, and other activity happens. Once that is implemented, the server will stay in sync with Zulip without needing to be restarted.

The backfill approach will likely remain useful for the initial load of historical messages, but the event API will take over for keeping the data current after startup.
