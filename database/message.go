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
	id_to_index map[int]int
	rows        []MessageRow
}

func NewMessageTable() *MessageTable {
	return &MessageTable{
		id_to_index: make(map[int]int),
		rows:        make([]MessageRow, 0),
	}
}

func (table *MessageTable) Put(message Message) int {
	id := message.MessageId

	index, ok := table.id_to_index[id]
	if ok {
		// for now we assume Messages never mutate
		// and we just got a repeat
		return index
	}

	new_index := len(table.rows)

	row := MessageRow{
		Index:        new_index,
		MessageId:    message.MessageId,
		Content:      message.Content,
		AddressIndex: message.AddressIndex,
		SenderIndex:  message.SenderIndex,
	}

	table.rows = append(table.rows, row)
	table.id_to_index[id] = new_index

	return new_index
}

func (table MessageTable) RowFromId(id int) *MessageRow {
	index, ok := table.id_to_index[id]
	if !ok {
		return nil
	}

	return &table.rows[index]
}
