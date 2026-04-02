package database

import "strings"

type TopicTable struct {
	Rows          []string
	name_to_index map[string]int
}

func NewTopicTable() *TopicTable {
	return &TopicTable{
		Rows:          make([]string, 0),
		name_to_index: make(map[string]int),
	}
}

func (table *TopicTable) Put(name string) int {
	index, ok := table.name_to_index[name]
	if ok {
		return index
	}

	new_index := len(table.Rows)

	table.Rows = append(table.Rows, strings.Clone(name))
	table.name_to_index[name] = new_index

	return new_index
}
