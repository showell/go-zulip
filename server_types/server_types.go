package servertypes

type ServerSubscription struct {
	StreamId int    `json:"stream_id"`
	Name     string `json:"name"`
}

type ServerMessage struct {
	Content        string `json:"content"`
	Id             int    `json:"id"`
	SenderFullName string `json:"sender_full_name"`
	SenderId       int    `json:"sender_id"`
	StreamId       int    `json:"stream_id"`
	Subject        string `json:"subject"`
	Type           string `json:"type"`
}
