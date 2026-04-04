package database

type Message struct {
	MessageId    int
	SenderIndex  int
	AddressIndex int
	Content      string
}

type MessageRow struct {
	Index        int
	MessageId    int
	SenderIndex  int
	AddressIndex int
	Content      string
}

type MessageTable struct {
	idToIndex map[int]int
	Rows      []MessageRow
}

func NewMessageTable() *MessageTable {
	return &MessageTable{
		idToIndex: make(map[int]int),
		Rows:      make([]MessageRow, 0),
	}
}

func (table *MessageTable) Put(message Message) int {
	id := message.MessageId

	index, ok := table.idToIndex[id]
	if ok {
		// for now we assume Messages never mutate
		// and we just got a repeat
		return index
	}

	newIndex := len(table.Rows)

	row := MessageRow{
		Index:        newIndex,
		MessageId:    message.MessageId,
		Content:      message.Content,
		AddressIndex: message.AddressIndex,
		SenderIndex:  message.SenderIndex,
	}

	table.Rows = append(table.Rows, row)
	table.idToIndex[id] = newIndex

	return newIndex
}

func (table MessageTable) RowFromId(id int) *MessageRow {
	index, ok := table.idToIndex[id]
	if !ok {
		return nil
	}

	return &table.Rows[index]
}
