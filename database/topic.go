package database

import "strings"

type TopicTable struct {
	Rows        []string
	nameToIndex map[string]int
}

func NewTopicTable() *TopicTable {
	return &TopicTable{
		Rows:        make([]string, 0),
		nameToIndex: make(map[string]int),
	}
}

func (table *TopicTable) Put(name string) int {
	index, ok := table.nameToIndex[name]
	if ok {
		return index
	}

	newIndex := len(table.Rows)

	table.Rows = append(table.Rows, strings.Clone(name))
	table.nameToIndex[name] = newIndex

	return newIndex
}
