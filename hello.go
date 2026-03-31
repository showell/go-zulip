package main

import "fmt"

type Channel struct {
    id int
    name string
}

type Database struct {
    channel_map map[int]string
}

func NewDatabase() *Database {
    return &Database{
        channel_map: make(map[int]string),
    }
}

func (db *Database) add_channel(channel Channel) {
    db.channel_map[channel.id] = channel.name
}

func (db Database) get_channel_name(id int) string {
    return db.channel_map[id]
}

func main() {
    channels := []Channel{
        {
            id: 101,
            name: "engineering",
        },
        {
            id: 102,
            name: "design",
        },
    }

    db := NewDatabase()

    for _, channel := range channels {
        db.add_channel(channel)
    }

    fmt.Println(db.get_channel_name(101))
}
