package main

import "testing"

func TestGMachine(t *testing.T) {
	t.Log(Run([]TSC{
		{
			Name: "f",
			Args: Vars("a"),
			Body: Var("a"),
		},
		{
			Name: "main",
			Args: nil,
			Body: Apps(Var("twice"), Var("f"), Num(42)),
		},
	}))

	//(define (take n l)
	//  (case l
	//    [(Nil) Nil]
	//    [(Cons x xs)
	//      (if (<= n 0)
	//        Nil
	//        (Cons x (take (- n 1) xs)))]))
}
