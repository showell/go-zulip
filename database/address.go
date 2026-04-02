package database

type AddressRow struct {
	ChannelIndex int
	TopicIndex   int
}

type AddressTable struct {
	row_to_index map[AddressRow]int
	Rows         []AddressRow
}

func NewAddressTable() *AddressTable {
	return &AddressTable{
		row_to_index: make(map[AddressRow]int),
		Rows:         make([]AddressRow, 0),
	}
}

func (table *AddressTable) Put(row AddressRow) int {
	index, ok := table.row_to_index[row]
	if ok {
		return index
	}

	new_index := len(table.Rows)

	table.Rows = append(table.Rows, row)
	table.row_to_index[row] = new_index

	return new_index
}
