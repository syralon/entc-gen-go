package service

func Trans[A, B any](a []A, fn func(A) B) []B {
	b := make([]B, 0, len(a))
	for _, v := range a {
		b = append(b, fn(v))
	}
	return b
}
