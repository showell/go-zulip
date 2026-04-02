package database

type AddressRow struct {
	ChannelIndex int
	TopicIndex   int
}

type AddressTable struct {
	row_to_index map[AddressRow]int
	rows         []AddressRow
}

func NewAddressTable() *AddressTable {
	return &AddressTable{
		row_to_index: make(map[AddressRow]int),
		rows:         make([]AddressRow, 0),
	}
}

func (table *AddressTable) Put(row AddressRow) int {
	index, ok := table.row_to_index[row]
	if ok {
		return index
	}

	new_index := len(table.rows)

	table.rows = append(table.rows, row)
	table.row_to_index[row] = new_index

	return new_index
}

func (table AddressTable) RowFromIndex(index int) *AddressRow {
	return &table.rows[index]
}
