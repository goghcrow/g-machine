package g_machine

import "testing"

const fpList = `
	(define (cons a b cc cn) (cc a b))
	(define (nil cc cn) cn)
	(define (hd list) (list K abort))
	(define (tl list) (list K1 abort))
	(define abort abort)
`

func TestGMachine(t *testing.T) {
	//	t.Log(Run(`
	//(define (f a) a)
	//(define main (twice f 42))
	//`))

	//	t.Log(Run(`
	//(define main (I 3))
	//`))

	//	t.Log(Run(`
	//(define id (S K K))
	//(define main (id 3))
	//`))

	//	t.Log(Run(`
	//(define id (S K K))
	//(define main (twice twice twice id 3))
	//`))

	//	t.Log(Run(`
	//(define main (twice (I I I) 3))
	//`))

	//cons a b cc cn = cc a b ;
	//nil cc cn = cn ;
	//hd list = list K abort ;
	//tl list = list K1 abort ;
	//abort = abort ;
	//infinite x = cons x (infinite x) ;
	//main = hd (tl (infinite 4))

	//t.Log(Run(fpList + `
	//(define (infinite x) (cons x (infinite x)))
	//(define main (hd (tl (infinite 4))))
	//`))

	//	t.Log(Run(`
	//(define main
	//	(let [(id1 (I I I))] (id1 id1 3)))
	//`))

	//t.Log(Run(`
	//(define (oct g x)
	//	(let [(h (twice g))]
	//		(let [(k (twice h))]
	//			(k (k x)))))
	//(define main (oct I 4))
	//`))

	// todo
	//t.Log(Run(`
	//(define (oct g x)
	//	(let [(h (twice g))
	//		  (k (twice h))]
	//		(k (k x))))
	//(define main (oct I 4))
	//`))

	// todo
	t.Log(Run(fpList + `
	(define (infinite x)
		(letrec [(xs (cons x xs))]
			xs))
	(define main (hd (tl (tl (infinite 4)))))
	`))

	//t.Log(Run(`
	//(define main (+ (* 4 5) (- 2 5)))
	//`))

	//	t.Log(Run(`
	//(define (inc x) (+ x 1))
	//(define main (twice twice inc 4))
	//	`))

	//t.Log(Run(fpList + `
	//(define (length xs) (xs length1 0))
	//(define (length1 x xs) (+ 1 (length xs)))
	//(define main (length (cons 3 (cons 3 (cons 3 nil)))))
	//`))

	//t.Log(Run(`
	//(define (fac n) (if (= n 0) 1 (* n (fac (- n 1)))))
	//(define main (fac 5))
	//`))

	//t.Log(Run(`
	//(define (gcd a b)
	//	(if (= a b)
	//		a
	//		(if (< a b)
	//			(gcd b a)
	//			(gcd b (- a b)))))
	//(define main (gcd 6 10))
	//`))

	//t.Log(Run(`
	//(define (add3 a b c) (+ (+ a b) c))
	//(define (nfib n)
	//	(if (= n 0)
	//		1
	//		(add3 1 (nfib (- n 1)) (nfib (- n 2)))))
	//(define main (nfib 4))
	//`))

	//t.Log(Run(`
	//(define (add3 a b c) (+ (+ a b) c))
	//(define (nfib n)
	//	(if (<= n 0)
	//		1
	//		(add3 1 (nfib (- n 1)) (nfib (- n 2)))))
	//(define main (nfib 4))
	//`))

	// todo
	//t.Log(Run(`
	//(define main (take 3 (sieve (from 2))))
	//(define (from n) (Cons n (from (+ n 1))))
	//(define (sieve xs)
	//	(case xs
	//		[(Nil) Nil]
	//		[(Cons p ps) (Cons p (sieve (filter (nonMultiple p) ps)))]))
	//(define (filter predicate xs)
	//	(case xs
	//		[(Nil) Nil]
	//		[(Cons p ps)
	//			(let [(rest (predicate ps))]
	//					(if (predicate p)
	//						(Cons p rest)
	//						rest))]))
	//(define (nonMultiple p n) (!= (* (/ n p) p) n))
	//(define (take n xs)
	//	(if (= n 0)
	//		Nil
	//		(case xs
	//			[(Nil) Nil]
	//			[(Cons p ps) (Cons p (take (- n 1) ps))])))
	//`))
}
