package database

import "strings"

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
	Rows        []IdNameRow
}

func NewIdNameTable() *IdNameTable {
	return &IdNameTable{
		id_to_index: make(map[int]int),
		Rows:        make([]IdNameRow, 0),
	}
}

func (table *IdNameTable) Put(id_name IdName) int {
	id := id_name.Id
	name := id_name.Name

	index, ok := table.id_to_index[id]
	if ok {
		return index
	}

	new_index := len(table.Rows)

	row := IdNameRow{
		Index: new_index,
		Id:    id,
		Name:  strings.Clone(name),
	}

	table.Rows = append(table.Rows, row)
	table.id_to_index[id] = new_index

	return new_index
}

func (table IdNameTable) RowFromId(id int) *IdNameRow {
	index, ok := table.id_to_index[id]
	if !ok {
		return nil
	}

	return &table.Rows[index]
}

func (table IdNameTable) GetOrMakeIndex(id int) int {
	index, ok := table.id_to_index[id]
	if ok {
		return index
	}

	return table.Put(IdName{Id: id, Name: "unknown"})
}

func (table IdNameTable) GetName(id int) string {
	row := table.RowFromId(id)
	if row == nil {
		return "unknown"
	}
	return row.Name
}
