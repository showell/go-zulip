package channel

type Channel struct {
	Id   int
	Name string
}

type ChannelRow struct {
	Index int
	Id    int
	Name  string
}

type ChannelTable struct {
	id_to_index map[int]int
	rows        []ChannelRow
}

func NewChannelTable() *ChannelTable {
	return &ChannelTable{
		id_to_index: make(map[int]int),
		rows:        make([]ChannelRow, 0),
	}
}

func (table *ChannelTable) Put(channel Channel) int {
	id := channel.Id
	name := channel.Name

	index, ok := table.id_to_index[id]
	if ok {
		return index
	}

	new_index := len(table.rows)

	row := ChannelRow{
		Index: new_index,
		Id:    id,
		Name:  name,
	}

	table.rows = append(table.rows, row)
	table.id_to_index[id] = new_index

	return new_index
}

func (table ChannelTable) RowFromId(id int) *ChannelRow {
	index, ok := table.id_to_index[id]
	if !ok {
		return nil
	}

	return &table.rows[index]
}
