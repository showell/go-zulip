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
	idToIndex map[int]int
	Rows      []IdNameRow
}

func NewIdNameTable() *IdNameTable {
	return &IdNameTable{
		idToIndex: make(map[int]int),
		Rows:      make([]IdNameRow, 0),
	}
}

func (table *IdNameTable) Put(idName IdName) int {
	id := idName.Id
	name := idName.Name

	index, ok := table.idToIndex[id]
	if ok {
		return index
	}

	newIndex := len(table.Rows)

	row := IdNameRow{
		Index: newIndex,
		Id:    id,
		Name:  strings.Clone(name),
	}

	table.Rows = append(table.Rows, row)
	table.idToIndex[id] = newIndex

	return newIndex
}

func (table IdNameTable) RowFromId(id int) *IdNameRow {
	index, ok := table.idToIndex[id]
	if !ok {
		return nil
	}

	return &table.Rows[index]
}

func (table IdNameTable) GetOrMakeIndex(id int) int {
	index, ok := table.idToIndex[id]
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
