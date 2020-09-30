package hash

import (
	"hash/fnv"
	"math"
	"sort"
)

type HashPicker interface {
	Pick(string) string
}

type consistentHash struct {
	cap          uint32
	sortedValues []sortValue
}

type sortValue struct {
	index uint32
	item  string
}

func NewConsistentHash(items []string, cap uint32) HashPicker {
	sortedValued := make([]sortValue, 0, cap)
	for _, item := range items {
		sortedValued = append(sortedValued, sortValue{
			index: hash(item, cap),
			item:  item,
		})
	}
	sort.Slice(sortedValued, func(i, j int) bool {
		return sortedValued[i].index < sortedValued[j].index
	})
	return &consistentHash{
		cap:          cap,
		sortedValues: sortedValued,
	}
}

func hash(s string, cap uint32) uint32 {
	aln := fnv.New32a()
	aln.Write([]byte(s))
	return aln.Sum32() % cap
}

func (c *consistentHash) Pick(s string) string {
	if len(c.sortedValues) == 0 {
		return ""
	}
	if len(c.sortedValues) == 1 {
		return c.sortedValues[0].item
	}
	lastIndex := len(c.sortedValues) - 1
	rightIndex := -1
	leftIndex := -1
	h := hash(s, c.cap)
	for i, item := range c.sortedValues {
		if item.index >= h {
			rightIndex = i
			break
		}
	}
	if rightIndex == -1 {
		rightIndex = 0
		leftIndex = lastIndex
	} else if rightIndex == 0 {
		leftIndex = lastIndex
	} else {
		leftIndex = rightIndex - 1
	}
	if math.Abs(float64(c.sortedValues[rightIndex].index-h)) < math.Abs(float64(h-c.sortedValues[leftIndex].index)) {
		return c.sortedValues[rightIndex].item
	} else {
		return c.sortedValues[leftIndex].item
	}
}
