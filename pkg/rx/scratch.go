package rx

func t() {
	var s Subscribable[int]

	Map(
		Map(s,
			func(value int) string {
				return ""
			}),
		func(value string) int {
			return 0
		})
}
