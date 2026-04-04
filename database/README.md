# database

An in-memory database that stores Zulip channels, topics, messages, and users.

The main entry point is `NewDatabase()`, which returns a `Database` struct containing the following tables:

- `ChannelTable` — channels (stream ID → name), stored as an `IdNameTable`
- `TopicTable` — topic names, deduplicated by string
- `AddressTable` — (channel, topic) pairs, each uniquely identifying a thread
- `MessageTable` — individual messages, keyed by message ID
- `UserTable` — senders (user ID → full name), stored as an `IdNameTable`

Two `OneToMany` indexes are maintained for fast lookups during rendering:
- `ChannelToAddress` — maps a channel to all its (channel, topic) addresses
- `AddressToMessage` — maps an address to all its messages

Data is loaded by calling `AddServerSubscription` and `AddServerMessage`, which accept the raw server types from the `server_types` package and handle all the indexing internally.
