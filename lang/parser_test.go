package lang

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	println(strings.Join(tokenize("(begin   (define \t\nr 10) (* pi (* r r)))"), "\n"))

	{
		pgrm := Parse(`
(define (take n l)
  (case l
    [(Nil) Nil]
    [(Cons x xs)
      (if (le n 0)
        Nil
        (Cons x (take (sub n 1) xs)))]))`)
		fmt.Println(pgrm)
	}

	{
		pgrm := Parse("(define (fibs) (Cons 0 (Cons 1 (zipWith add fibs (tail fibs)))))")
		fmt.Println(pgrm)
	}

	{
		pgrm := Parse(`
(define a (let ([a 1] [b (+ a 1)]) (+ a b)))
`)
		fmt.Println(pgrm)
	}

	fmt.Println(strings.Join(SliceMap(Parse(PreludeDefs), func(t TSC) string { return t.String() }), "\n"))
}

var PreludeDefs = `
(define (i x) x)
(define (k x y) x)
(define (k1 x y) y)
(define (a f g x) (f x (g x)))
(define (compose f g x) (f (g x)))
(define (twice f) (compose f f))
`
