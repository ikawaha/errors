package linker

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestLink(t *testing.T) {
	t.Run("append: nil, empty", func(t *testing.T) {
		e := Append(nil)
		if e != nil {
			t.Errorf("want nil, but %+v", e)
		}
	})
	t.Run("append: nil, nil", func(t *testing.T) {
		e := Append(nil, nil)
		if e != nil {
			t.Errorf("want nil, but %+v", e)
		}
	})
	t.Run("append: err, nil", func(t *testing.T) {
		lhs := errors.New("lhs")
		e := Append(lhs, nil)
		if e.Error() != lhs.Error() {
			t.Errorf("want: %v, got %+v", lhs, e)
		}
		if !errors.Is(e, lhs) {
			t.Error("want: errors.Is(e, lhs)=true, but false")
		}
	})
	t.Run("append: nil, err", func(t *testing.T) {
		rhs := errors.New("lhs")
		e := Append(nil, rhs)
		if e.Error() != rhs.Error() {
			t.Errorf("want: %v, got %+v", rhs, e)
		}
		if !errors.Is(e, rhs) {
			t.Error("want: errors.Is(e, rhs)=true, but false")
		}
	})
	t.Run("append: err, err", func(t *testing.T) {
		lhs := errors.New("lhs")
		rhs := errors.New("rhs")
		e := Append(lhs, rhs)

		fmt.Println("errors.Is(e, lhs)=", errors.Is(e, lhs))
		fmt.Println("errors.Is(e, rhs)=", errors.Is(e, rhs))
		fmt.Println("yield errors=", Yield(e))
		if e.Error() != lhs.Error() {
			t.Errorf("want: %v, got %+v", lhs, e)
		}
		if !errors.Is(e, lhs) {
			t.Error("want: errors.Is(e, lhs)=true, but false")
		}
		if !errors.Is(e, rhs) {
			t.Error("want: errors.Is(e, rhs)=true, but false")
		}
	})
	t.Run("append: err, err list", func(t *testing.T) {
		lhs := errors.New("lhs")
		var rhs []error
		for i := 0; i < 10; i++ {
			rhs = append(rhs, errors.New("rhs_"+strconv.Itoa(i)))
		}
		e := Append(lhs, rhs...)
		if e.Error() != lhs.Error() {
			t.Errorf("want: %v, got %+v", lhs, e)
		}
		if !errors.Is(e, lhs) {
			t.Errorf("want: errors.Is(e, lhs)=true, but false")
		}
		for i := 0; i < 10; i++ {
			if !errors.Is(e, rhs[i]) {
				t.Errorf("want: errors.Is(e, rhs[%d])=true, but false", i)
			}
		}
	})
	t.Run("append: nil, err list", func(t *testing.T) {
		var rhs []error
		for i := 0; i < 10; i++ {
			rhs = append(rhs, errors.New("rhs_"+strconv.Itoa(i)))
		}
		e := Append(nil, rhs...)
		if e.Error() != rhs[0].Error() {
			t.Errorf("want: %v, got %+v", rhs[0], e)
		}
		for i := 0; i < 10; i++ {
			if !errors.Is(e, rhs[i]) {
				t.Errorf("want: errors.Is(e, rhs[%d])=true, but false", i)
			}
		}
	})
}

func TestYield(t *testing.T) {
	t.Run("yield: append[nil]", func(t *testing.T) {
		e := Append(nil)
		got := Yield(e)
		if got != nil {
			t.Errorf("want: nil, got: %v", got)
		}
	})
	t.Run("yield: append[nil, nil]", func(t *testing.T) {
		e := Append(nil, nil)
		got := Yield(e)
		if got != nil {
			t.Errorf("want: nil, got: %v", got)
		}
	})
	t.Run("yield: append[err, nil]", func(t *testing.T) {
		lhs := errors.New("lhs")
		e := Append(lhs, nil)
		got := Yield(e)
		want := []error{lhs}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %+v, got: %+v", want, got)
		}
	})
	t.Run("yield: append[err, err]", func(t *testing.T) {
		lhs := errors.New("lhs")
		rhs := errors.New("rhs")
		e := Append(lhs, rhs)
		got := Yield(e)
		want := []error{lhs, rhs}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %+v, got: %+v", want, got)
		}
	})
	t.Run("yield: append[err, err list]", func(t *testing.T) {
		lhs := errors.New("lhs")
		var rhs []error
		for i := 0; i < 10; i++ {
			rhs = append(rhs, errors.New("rhs_"+strconv.Itoa(i)))
		}
		e := Append(lhs, rhs...)
		got := Yield(e)
		want := append([]error{lhs}, rhs...)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %+v, got: %+v", want, got)
		}
	})
	t.Run("yield: append[append[err_list], append[err_list]]", func(t *testing.T) {
		var lhs []error
		for i := 0; i < 10; i++ {
			lhs = append(lhs, errors.New("lhs_"+strconv.Itoa(i)))
		}
		var rhs []error
		for i := 0; i < 10; i++ {
			rhs = append(rhs, errors.New("rhs_"+strconv.Itoa(i)))
		}
		e := Append(Append(lhs[0], lhs[1:]...), Append(rhs[0], rhs[1:]...))
		got := Yield(e)
		want := append(lhs, rhs...)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %+v, got: %+v", want, got)
		}
	})
	t.Run("yield: embedding", func(t *testing.T) {
		var lhs, want []error
		for i := 0; i < 3; i++ {
			var es []error
			for j := 0; j < 3; j++ {
				es = append(es, fmt.Errorf("lhs%d_%d", i, j))
			}
			lhs = append(lhs, Append(es[0], es[1:]...))
			want = append(want, es...)
		}
		var rhs []error
		for i := 0; i < 3; i++ {
			rhs = append(rhs, errors.New("rhs_"+strconv.Itoa(i)))
		}
		want = append(want, rhs...)

		e := Append(Append(lhs[0], lhs[1:]...), Append(rhs[0], rhs[1:]...))
		got := Yield(e)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %+v, got: %+v", want, got)
		}
	})
}
