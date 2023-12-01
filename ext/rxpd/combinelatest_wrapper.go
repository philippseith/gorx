package rxpd

import "time"

func CombineLatest2[S1 any, S2 any, T any](combine func(S1, S2) T,
	s1 Subscribable[S1], s2 Subscribable[S2]) Property[T] {
	return CombineLatest[T](func(next ...any) T {
		return combine(next[0].(S1), next[1].(S2))
	}, ToAny(s1), ToAny(s2))
}

func CombineLatest3[S1 any, S2 any, S3 any, T any](combine func(S1, S2, S3) T,
	s1 Subscribable[S1], s2 Subscribable[S2], s3 Subscribable[S3]) Property[T] {
	return CombineLatest[T](func(next ...any) T {
		return combine(next[0].(S1), next[1].(S2), next[2].(S3))
	}, ToAny(s1), ToAny(s2), ToAny(s3))
}

func CombineLatest4[S1 any, S2 any, S3 any, S4 any, T any](combine func(S1, S2, S3, S4) T,
	s1 Subscribable[S1], s2 Subscribable[S2], s3 Subscribable[S3], s4 Subscribable[S4]) Property[T] {
	return CombineLatest[T](func(next ...any) T {
		return combine(next[0].(S1), next[1].(S2), next[2].(S3), next[3].(S4))
	}, ToAny(s1), ToAny(s2), ToAny(s3), ToAny(s4))
}

func CombineLatest5[S1 any, S2 any, S3 any, S4 any, S5 any, T any](combine func(S1, S2, S3, S4, S5) T,
	s1 Subscribable[S1], s2 Subscribable[S2], s3 Subscribable[S3], s4 Subscribable[S4],
	s5 Subscribable[S5]) Property[T] {
	return CombineLatest[T](func(next ...any) T {
		return combine(next[0].(S1), next[1].(S2), next[2].(S3), next[3].(S4), next[4].(S5))
	}, ToAny(s1), ToAny(s2), ToAny(s3), ToAny(s4), ToAny(s5))
}

func CombineLatest6[S1 any, S2 any, S3 any, S4 any, S5 any, S6 any, T any](
	combine func(S1, S2, S3, S4, S5, S6) T,
	s1 Subscribable[S1], s2 Subscribable[S2], s3 Subscribable[S3], s4 Subscribable[S4],
	s5 Subscribable[S5], s6 Subscribable[S6]) Property[T] {
	return CombineLatest[T](func(next ...any) T {
		return combine(next[0].(S1), next[1].(S2), next[2].(S3), next[3].(S4), next[4].(S5), next[5].(S6))
	}, ToAny(s1), ToAny(s2), ToAny(s3), ToAny(s4), ToAny(s5), ToAny(s6))
}

func CombineLatest7[S1 any, S2 any, S3 any, S4 any, S5 any, S6 any, S7 any, T any](
	combine func(S1, S2, S3, S4, S5, S6, S7) T,
	s1 Subscribable[S1], s2 Subscribable[S2], s3 Subscribable[S3], s4 Subscribable[S4],
	s5 Subscribable[S5], s6 Subscribable[S6], s7 Subscribable[S7]) Property[T] {
	return CombineLatest[T](func(next ...any) T {
		return combine(next[0].(S1), next[1].(S2), next[2].(S3), next[3].(S4), next[4].(S5), next[5].(S6),
			next[6].(S7))
	}, ToAny(s1), ToAny(s2), ToAny(s3), ToAny(s4), ToAny(s5), ToAny(s6), ToAny(s7))
}

func CombineLatest8[S1 any, S2 any, S3 any, S4 any, S5 any, S6 any, S7 any, S8 any, T any](
	combine func(S1, S2, S3, S4, S5, S6, S7, S8) T,
	s1 Subscribable[S1], s2 Subscribable[S2], s3 Subscribable[S3], s4 Subscribable[S4],
	s5 Subscribable[S5], s6 Subscribable[S6], s7 Subscribable[S7], s8 Subscribable[S8]) Property[T] {
	return CombineLatest[T](func(next ...any) T {
		return combine(next[0].(S1), next[1].(S2), next[2].(S3), next[3].(S4), next[4].(S5), next[5].(S6),
			next[6].(S7), next[7].(S8))
	}, ToAny(s1), ToAny(s2), ToAny(s3), ToAny(s4), ToAny(s5), ToAny(s6), ToAny(s7), ToAny(s8))
}

func CombineLatest9[S1 any, S2 any, S3 any, S4 any, S5 any, S6 any, S7 any, S8 any, S9 any, T any](
	combine func(S1, S2, S3, S4, S5, S6, S7, S8, S9) T,
	s1 Subscribable[S1], s2 Subscribable[S2], s3 Subscribable[S3], s4 Subscribable[S4],
	s5 Subscribable[S5], s6 Subscribable[S6], s7 Subscribable[S7], s8 Subscribable[S8],
	s9 Subscribable[S9]) Property[T] {
	return CombineLatest[T](func(next ...any) T {
		return combine(next[0].(S1), next[1].(S2), next[2].(S3), next[3].(S4), next[4].(S5), next[5].(S6),
			next[6].(S7), next[7].(S8), next[8].(S9))
	}, ToAny(s1), ToAny(s2), ToAny(s3), ToAny(s4), ToAny(s5), ToAny(s6), ToAny(s7), ToAny(s8), ToAny(s9))
}

func CombineLatestToSlice[T any](sst []Subscribable[T]) Property[[]T] {
	ssa := make([]Subscribable[any], 0, len(sst))
	for _, st := range sst {
		ssa = append(ssa, ToAny(st))
	}
	return CombineLatest[[]T](func(next ...any) []T {
		result := make([]T, 0, len(next))
		for _, n := range next {
			result = append(result, n.(T))
		}
		return result
	}, ssa...).With(WithSetInterval[[]T](func(interval time.Duration) {
		for _, st := range sst {
			st.SetInterval(interval)
		}
	}))
}
