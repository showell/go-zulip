package database

import (
	"maps"
	"slices"
)

type IntSet = map[int]struct{}

type OneToMany struct {
	list_of_sets []IntSet
}

func NewOneToMany() *OneToMany {
	return &OneToMany{
		list_of_sets: make([]IntSet, 0),
	}
}

func (self *OneToMany) Update(one_index int, many_index int) {
	for one_index >= len(self.list_of_sets) {
		int_set := make(IntSet)
		self.list_of_sets = append(self.list_of_sets, int_set)
	}

	many_index_set := self.list_of_sets[one_index]
	many_index_set[many_index] = struct{}{}
	self.list_of_sets[one_index] = many_index_set
}

func (self OneToMany) GetManyIndexesInRandomOrder(one_index int) []int {
	int_set := self.list_of_sets[one_index]

	if len(int_set) == 0 {
		return []int{}
	}

	return slices.Collect(maps.Keys(int_set))
}

func (self OneToMany) Count(one_index int) int {
	return len(self.list_of_sets[one_index])
}
