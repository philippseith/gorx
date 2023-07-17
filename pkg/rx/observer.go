package rx

type Observer[T any] interface {
	Next(value T)
	Error(err error)
	Complete()
}
