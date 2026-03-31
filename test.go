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
    }

    fmt.Println(db.GetChannelName(101))
}
