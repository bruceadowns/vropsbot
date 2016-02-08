package common

import (
	"testing"
)

var bmtests = []struct {
	key, value string
}{
	{
		key:   "foo",
		value: "bar",
	},
	{
		key:   "ham",
		value: "eggs",
	}}

func TestBiMap(t *testing.T) {
	bm := NewBiMap()

	for _, test := range bmtests {
		bm.Put(test.key, test.value)
	}

	for _, test := range bmtests {
		v, e := bm.GetByKey(test.key)
		if !e {
			t.Errorf("%s was not found. Expected %s.", test.key, test.value)
		}
		if v != test.value {
			t.Errorf("%s Expected %s. Actual %s", test.key, test.value, v)
		}

		k, e := bm.GetByValue(test.value)
		if !e {
			t.Errorf("%s was not found. Expected %s.", test.value, test.key)
		}
		if k != test.key {
			t.Errorf("%s Expected %s. Actual %s", test.value, test.key, k)
		}
	}
}
