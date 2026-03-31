package database

type IdName struct {
	Id   int
	Name string
}

type IdNameRow struct {
	Index int
	Id    int
	Name  string
}

type IdNameTable struct {
	id_to_index map[int]int
	rows        []IdNameRow
}

func NewIdNameTable() *IdNameTable {
	return &IdNameTable{
		id_to_index: make(map[int]int),
		rows:        make([]IdNameRow, 0),
	}
}

func (table *IdNameTable) Put(id_name IdName) int {
	id := id_name.Id
	name := id_name.Name

	index, ok := table.id_to_index[id]
	if ok {
		return index
	}

	new_index := len(table.rows)

	row := IdNameRow{
		Index: new_index,
		Id:    id,
		Name:  name,
	}

	table.rows = append(table.rows, row)
	table.id_to_index[id] = new_index

	return new_index
}

func (table IdNameTable) RowFromId(id int) *IdNameRow {
	index, ok := table.id_to_index[id]
	if !ok {
		return nil
	}

	return &table.rows[index]
}

func (table IdNameTable) GetName(id int) string {
	row := table.RowFromId(id)
	if row == nil {
		return "unknown"
	}
	return row.Name
}
