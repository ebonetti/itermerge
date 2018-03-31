// Package itermerge is a package that provides primitives for an heap of iterators.
// It's a generic package when used in conjunction with github.com/taylorchu/generic
package itermerge

import "testing"

const (
	MIN = 0
	MAX = 10000
)

func Test(t *testing.T) {
	gc := generator(MIN, MAX)
	iterMerge := New(chan2func(gc))

	if v, ok := iterMerge.Peek(); v.(myType) != myType(MIN) || !ok {
		t.Errorf("Peek should return %v,%v, instead it returns %v,%v", MIN, true, v, ok)
	}

	vtest := myType(MIN)
	for v, ok := iterMerge.Next(); ok; v, ok = iterMerge.Next() {
		if vtest != v {
			t.Errorf("Next should return the value %v, instead it returns %v", vtest, v)
		}
		vtest++
		iterMerge.Push(chan2func(gc)) //add another worker/path for testing purposes
	}

	if vtest != MAX+1 {
		t.Errorf("vtest at the end is %v, but it should be %v", vtest, MAX+1)
	}
}

func chan2func(c <-chan Type) func() (Type, bool) {
	return func() (t Type, ok bool) {
		t, ok = <-c
		return
	}
}

func generator(min, max myType) <-chan Type {
	result := make(chan Type, 100)

	go func() {
		defer close(result)
		for i := min; i <= max; i++ {
			result <- i
		}
	}()

	return result
}

type myType int

func (x myType) Less(y interface{}) bool {
	return x < y.(myType)
}
