package lazy

import "fmt"

type ShowDict[T any] struct {
	Show func(T) string
}

var instanceShowDictForI32 = ShowDict[int32]{
	Show: func(i int32) string {
		return fmt.Sprintf("%d", i)
	},
}

type Point[T ~int | int64] struct {
	X, Y int64
}

var instanceShowDictForPoint = ShowDict[Point[int64]]{
	Show: func(p Point[int64]) string {
		return fmt.Sprintf("Point(%d, %d)", p.X, p.Y)
	},
}

func Show[T any](dict ShowDict[T], a T) string {
	return dict.Show(a)
}
