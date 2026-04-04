# address

This package handles URL routing for the web server. It parses a URL path string into a typed Go value that callers can pattern-match on to decide what to render.

## Address types

There are four address types, all implementing the `Address` interface:

- `ChannelsAddress` — matches `/channels`
- `TopicsAddress` — matches `/topics/{channelId}`, carries the channel ID
- `MessagesAddress` — matches `/messages/{addressIndex}`, carries the address index
- `NadaAddress` — the fallback for unrecognized paths

## Code pattern

`GetAddress` is the main entry point. It uses simple string matching for exact paths (`/channels`) and compiled regexes for paths with numeric parameters (`/topics/(\d+)`, `/messages/(\d+)`). Each successful match returns a concrete struct value populated with the parsed integer.

The caller (`html.Html`) does a type switch on the returned `Address` value to dispatch to the appropriate HTML renderer, passing along any numeric fields from the struct.

## WritePath

Each address type implements `WritePath`, which serializes the address back to a URL path string. This is used in the HTML layer when building anchor tags — for example, the channels page constructs `TopicsAddress` values and calls `WritePath` to write the `href` directly into the output stream, and similarly the topics page does the same with `MessagesAddress`. Keeping path formatting in one place avoids duplicating URL construction logic across the HTML renderers.
