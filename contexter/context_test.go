package contexter_test

import (
	"errors"
	"testing"

	"github.com/ikawaha/errors/contexter"
)

func TestWithValue(t *testing.T) {
	err := errors.New("")
	err = contexter.WithValue(err, "key1", "value1")
	err = contexter.WithValue(err, "key2", "value2")
	err = contexter.WithValue(err, "key3", "value3")
	err = contexter.WithValue(err, "key1", "overwritten")
	err = contexter.WithValue(err, 123, -123)

	testdata := []struct {
		name    string
		key     any
		wantOk  bool
		wantVal any
	}{
		{name: "ok: key2", key: "key2", wantOk: true, wantVal: "value2"},
		{name: "ok: key3", key: "key3", wantOk: true, wantVal: "value3"},
		{name: "ok: key1 (overwritten)", key: "key1", wantOk: true, wantVal: "overwritten"},
		{name: "ok: 123", key: 123, wantOk: true, wantVal: -123},
		{name: "ng (empty key)", key: nil, wantOk: false, wantVal: nil},
		{name: "ng", key: "hello", wantOk: false, wantVal: nil},
	}
	for _, d := range testdata {
		t.Run(d.name, func(t *testing.T) {
			got, ok := contexter.Value(err, d.key)
			if ok != d.wantOk {
				t.Errorf("want: %t, got: %t", d.wantOk, ok)
			}
			if got != d.wantVal {
				t.Errorf("want: %v, got: %v", d.wantVal, got)
			}
		})
	}
}
