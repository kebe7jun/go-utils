package flatten

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Separator struct {
	Before string
	After  string
}

type Flattener struct {
	data    string
	dictSep Separator
	listSep Separator
}

var defaultFlattener = Flattener{
	dictSep: Separator{
		Before: ".",
	},
	listSep: Separator{
		Before: "[",
		After:  "]",
	},
}

func NewFlattener(data string) Flattener {
	return Flattener{
		data:    data,
		dictSep: defaultFlattener.dictSep,
		listSep: defaultFlattener.listSep,
	}
}

func NewFlattenerWithSep(data string, dictSep Separator, listSep Separator) Flattener {
	return Flattener{
		data:    data,
		dictSep: dictSep,
		listSep: listSep,
	}
}

func (f *Flattener) Flatten() (map[string]interface{}, error) {
	res := make(map[string]interface{})
	if strings.HasPrefix(strings.TrimSpace(f.data), "[") {
		d := make([]interface{}, 0)
		err := json.Unmarshal([]byte(f.data), &d)
		if err != nil {
			return nil, err
		}
		err = f.flatList("", d, res)
		if err != nil {
			return nil, err
		}
	} else {
		d := make(map[string]interface{}, 0)
		err := json.Unmarshal([]byte(f.data), &d)
		if err != nil {
			return nil, err
		}
		err = f.flatMap("", d, res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (f *Flattener) flatMap(prefix string, m, res map[string]interface{}) error {
	for k, d := range m {
		p := f.makeKey(prefix, k, false)
		switch d.(type) {
		case string, bool, float64, nil:
			res[p] = d
		case []interface{}:
			err := f.flatList(p, d.([]interface{}), res)
			if err != nil {
				return err
			}
		case map[string]interface{}:
			err := f.flatMap(p, d.(map[string]interface{}), res)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported type: %t", d)
		}
	}
	return nil
}

func (f *Flattener) flatList(prefix string, l []interface{}, res map[string]interface{}) error {
	for i, inf := range l {
		p := f.makeKey(prefix, fmt.Sprintf("%d", i), true)
		switch inf.(type) {
		case string, bool, float64, nil:
			res[p] = inf
		case map[string]interface{}:
			err := f.flatMap(p, inf.(map[string]interface{}), res)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported type: %t", inf)
		}
	}
	return nil
}

func (f *Flattener) makeKey(prefix, key string, isListItem bool) string {
	sep := Separator{}
	if isListItem {
		sep = f.listSep
	} else {
		sep = f.dictSep
	}
	return fmt.Sprintf("%s%s%s%s", prefix, sep.Before, key, sep.After)
}
