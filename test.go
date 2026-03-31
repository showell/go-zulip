package main

import "fmt"
import "zulip-go/database"

func main() {
    channels := []database.Channel{
        {
            Id: 101,
            Name: "engineering",
        },
        {
            Id: 102,
            Name: "design",
        },
    }

    db := database.NewDatabase()

    for _, channel := range channels {
        db.AddChannel(channel)
    }

    fmt.Println(db.GetChannelName(101))
}
