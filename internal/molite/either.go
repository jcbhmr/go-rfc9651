package molite

type Either[A, B any] struct {
	Which int
	A     A
	B     B
}

type Either3[A, B, C any] struct {
	Which int
	A     A
	B     B
	C     C
}
