package server_types

type ServerSubscription struct {
	Stream_id int
	Name      string
}

type ServerMessage struct {
	Content          string
	Id               int
	Sender_full_name string
	Sender_id        int
	Subject          string
	Stream_id        int
}
