package flatten

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlatten(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		wannaRes map[string]interface{}
		wannaErr bool
	}{
		{
			"normal",
			//`{"a": 1, "b": true, "c": null, "d": [1], "e": {"a": 1}}`,
			`{"a": 1, "b": true}`,
			map[string]interface{}{
				".a": float64(1),
				".b": true,
			},
			false,
		},
		{
			"all types",
			`{"a": 1, "b": true, "c": null, "d": [1], "e": {"a": 1}}`,
			map[string]interface{}{
				".a":    float64(1),
				".b":    true,
				".c":    nil,
				".d[0]": float64(1),
				".e.a":  float64(1),
			},
			false,
		},
		{
			"list",
			`[1,null,false,"1"]`,
			map[string]interface{}{
				"[0]": float64(1),
				"[1]": nil,
				"[2]": false,
				"[3]": "1",
			},
			false,
		},
		{
			"deep",
			`{"a": {"b": {"c": 1, "d": [{"e": false}]}}}`,
			map[string]interface{}{
				".a.b.c":      float64(1),
				".a.b.d[0].e": false,
			},
			false,
		},
		{
			"json format error",
			`{xxx}`,
			nil,
			true,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d---%s", i, c.name), func(t *testing.T) {
			r := NewFlattener(c.in)
			res, err := r.Flatten()
			if c.wannaErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, c.wannaRes, res)
		})
	}
}

func TestFlattener_FlattenCustomSep(t *testing.T) {
	r := NewFlattenerWithSep(`{"a": 1, "b": [true]}`, Separator{
		Before: "<",
		After:  ">",
	}, Separator{
		Before: "{",
		After:  "}",
	})
	res, err := r.Flatten()
	assert.NoError(t, err)
	assert.Equal(t, res, map[string]interface{}{
		"<a>":    float64(1),
		"<b>{0}": true,
	})
}
