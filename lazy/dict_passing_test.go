package lazy

import "testing"

func TestDictPassing(t *testing.T) {
	println(Show[int32](instanceShowDictForI32, 42))

	println(Show[Point[int64]](instanceShowDictForPoint, Point[int64]{42, 100}))
}
