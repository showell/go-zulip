package database

type Channel struct {
    Id int
    Name string
}

type Database struct {
    channel_map map[int]string
}

func NewDatabase() *Database {
    return &Database{
        channel_map: make(map[int]string),
    }
}

func (db *Database) AddChannel(channel Channel) {
    db.channel_map[channel.Id] = channel.Name
}

func (db Database) GetChannelName(id int) string {
    return db.channel_map[id]
}

