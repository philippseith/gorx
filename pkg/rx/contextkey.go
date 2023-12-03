package rx

// ContextKey should be used as key for context.WithValue
type ContextKey string

func (c ContextKey) String() string {
	return "rx: " + string(c)
}

var (
	// ContextKeyCreate should be used to enrich contexts passed to observers of the Create function
	ContextKeyCreate = ContextKey("Create")
	// ContextKeyDefer should be used to enrich contexts passed to observers of the Defer function
	ContextKeyDefer = ContextKey("Defer")
)
