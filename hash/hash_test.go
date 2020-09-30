package hash

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConsistentHash_Pick(t *testing.T) {
	cases := []struct {
		name  string
		cap   uint32
		items []string
		in    string
		out   string
	}{
		{
			"normal",
			128,
			[]string{"1", "2", "3"},
			"222",
			"3",
		},
		{
			"normal-4",
			128,
			[]string{"1", "2", "3", "4"},
			"222",
			"3",
		},
		{
			"normal-5",
			128,
			[]string{"1", "2", "3", "5"},
			"222",
			"3",
		},
		{
			"normal-8",
			128,
			[]string{"1", "2", "3", "5", "6", "7", "8"},
			"222",
			"3",
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d---%s", i, c.name), func(t *testing.T) {
			hash := NewConsistentHash(c.items, c.cap)
			assert.Equal(t, hash.Pick(c.in), c.out)
		})
	}
}
