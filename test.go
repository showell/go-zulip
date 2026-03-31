package main

import "fmt"
import "zulip-go/database"
import "zulip-go/server_types"

type ServerSubscription = server_types.ServerSubscription

func main() {
    subs := []ServerSubscription{
        {
            StreamId: 101,
            Name: "engineering",
        },
        {
            StreamId: 102,
            Name: "design",
        },
    }

    db := database.NewDatabase()

    for _, sub := range subs {
        db.AddServerSubscription(sub)
        // test idempotency
        index := db.AddServerSubscription(sub)
        fmt.Println(index)
    }

    fmt.Println(db.GetChannelName(101))
    fmt.Println(db.GetChannelName(102))
}
