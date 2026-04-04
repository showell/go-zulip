package database

type AddressRow struct {
	ChannelIndex int
	TopicIndex   int
}

type AddressTable struct {
	rowToIndex map[AddressRow]int
	Rows       []AddressRow
}

func NewAddressTable() *AddressTable {
	return &AddressTable{
		rowToIndex: make(map[AddressRow]int),
		Rows:       make([]AddressRow, 0),
	}
}

func (table *AddressTable) Put(row AddressRow) int {
	index, ok := table.rowToIndex[row]
	if ok {
		return index
	}

	newIndex := len(table.Rows)

	table.Rows = append(table.Rows, row)
	table.rowToIndex[row] = newIndex

	return newIndex
}
