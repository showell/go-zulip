package topic

import "strings"

type TopicTable struct {
	rows          []string
	name_to_index map[string]int
}

func NewTopicTable() *TopicTable {
	return &TopicTable{
		rows:          make([]string, 0),
		name_to_index: make(map[string]int),
	}
}

func (table *TopicTable) Put(name string) int {
	index, ok := table.name_to_index[name]
	if ok {
		return index
	}

	new_index := len(table.rows)

	table.rows = append(table.rows, strings.Clone(name))
	table.name_to_index[name] = new_index

	return new_index
}

func (table TopicTable) NameFromIndex(index int) string {
	return table.rows[index]
}
