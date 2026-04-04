package database

import (
	"maps"
	"slices"
)

type IntSet = map[int]struct{}

type OneToMany struct {
	listOfSets []IntSet
}

func NewOneToMany() *OneToMany {
	return &OneToMany{
		listOfSets: make([]IntSet, 0),
	}
}

func (self *OneToMany) Update(oneIndex int, manyIndex int) {
	for oneIndex >= len(self.listOfSets) {
		intSet := make(IntSet)
		self.listOfSets = append(self.listOfSets, intSet)
	}

	manyIndexSet := self.listOfSets[oneIndex]
	manyIndexSet[manyIndex] = struct{}{}
	self.listOfSets[oneIndex] = manyIndexSet
}

func (self OneToMany) GetManyIndexesInRandomOrder(oneIndex int) []int {
	if oneIndex >= len(self.listOfSets) {
		return []int{}
	}

	intSet := self.listOfSets[oneIndex]

	if len(intSet) == 0 {
		return []int{}
	}

	return slices.Collect(maps.Keys(intSet))
}

func (self OneToMany) Count(oneIndex int) int {
	if oneIndex >= len(self.listOfSets) {
		return 0
	}
	return len(self.listOfSets[oneIndex])
}
