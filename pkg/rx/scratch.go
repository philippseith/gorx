package rx

func t() {
	var s Subscribable[int]

	Pipe(
		Pipe(s,
			func(value int) string {
				return ""
			}),
		func(value string) int {
			return 0
		})
}
