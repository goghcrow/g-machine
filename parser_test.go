package main

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
}
