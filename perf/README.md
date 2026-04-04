# perf

Performance testing harness for the HTML rendering code.

## Why it's fast

The HTML rendering is fast because the underlying data structures are designed so that rendering never has to do any hash lookups or pointer chasing through complex structures:

- **Everything is indexed by integer** — channels, topics, addresses, users, and messages are all stored in flat slices and looked up by integer index. Array access by index is O(1) and cache-friendly. There are no map lookups in the hot rendering path.
- **`OneToMany` uses integer indexes on both sides** — `ChannelToAddress` and `AddressToMessage` map an integer to a set of integers. During rendering you just grab a slice of ints and iterate — no string keys, no hashing.
- **Strings are stored once** — channel names, topic names, sender names, and message content are each stored in one place and referenced by index everywhere else. The renderer just does `table.Rows[index]` to get the string, which is a single array dereference.
- **`AddressRow` is a struct of two ints** — small, flat, and cheap to compare.
- **The database is built once and then read-only during rendering** — no locking, no synchronization overhead, no defensive copying.
- **Direct writes to `io.StringWriter`** — no intermediate buffers or string concatenation. Each piece is written straight to the destination with no temporary allocations.
- **No template engine** — Go's `html/template` does reflection and allocation on every render. The rendering here is just a sequence of direct function calls that the compiler can inline aggressively.

The overall effect is that rendering a page is essentially just array indexing and string writing — about as fast as you can get in a garbage-collected language.

## Benchmark results

`buildBigDb()` generates a synthetic database of 20 channels × 20 topics × 1000 message batches (400,000 messages total).

- **`channels()`** — renders the channels page 10 million times, producing ~33 billion characters of HTML
- **`topicsAndMessages()`** — renders topics and messages pages until 10 billion characters are produced (~320 outer loops across all 20 channels and their topics)

## Running

Uncomment the desired function in `main()` and run:

```bash
go run ./perf/
```
