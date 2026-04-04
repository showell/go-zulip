package servertypes

type ServerSubscription struct {
	StreamId int
	Name     string
}

type ServerMessage struct {
	Content        string
	Id             int
	SenderFullName string
	SenderId       int
	StreamId       int
	Subject        string
}
