package lazy

import (
	"reflect"
	"runtime/debug"
	"testing"
)

func assert(t *testing.T, a bool) {
	if !a {
		t.Fail()
	}
}
func assertEqual(t *testing.T, a any, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("expected %v but got %v", b, a)
		println(string(debug.Stack()))
	}
}

func TestList(t *testing.T) {
	assert(t, len(ListTake[int](ListFrom([]int(nil)), 100)) == 0)

	assertEqual(t, ListTake[int](ListFrom([]int{1, 2, 3}), -1), []int{})
	assertEqual(t, ListTake[int](ListFrom([]int{1, 2, 3}), 0), []int{})
	assertEqual(t, ListTake[int](ListFrom([]int{1, 2, 3}), 1), []int{1})
	assertEqual(t, ListTake[int](ListFrom([]int{1, 2, 3}), 2), []int{1, 2})
	assertEqual(t, ListTake[int](ListFrom([]int{1, 2, 3}), 3), []int{1, 2, 3})
	assertEqual(t, ListTake[int](ListFrom([]int{1, 2, 3}), 100), []int{1, 2, 3})

	assertEqual(t, ListTo[int](ListMap(ListFrom([]int{1, 2, 3}), func(x int) int {
		return x + 1
	})), []int{2, 3, 4})
}
