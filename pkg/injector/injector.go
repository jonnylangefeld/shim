package injector

// Run is the main utility function this package offers.
// It allows to run with replacements. A `Replacer` can be created via
// `Replace(&original).With(replacement)`.
//
// Example:
func Run(f func(), replacements ...Replacer) {
	for _, i := range replacements {
		i.inject()
	}

	f()

	for _, i := range replacements {
		i.revert()
	}
}

// Replacer is the interface that handles the replacement before and after
// the function execution.
type Replacer interface {
	inject()
	revert()
}

// replacer is the generic type implementing the interface. It stores the
// actual underlying original and replacement.
type replacer[T any] struct {
	original               *T
	replacement, origCache T
}

// Replace is the entry function to instantiate a replacer.
// It takes the pointer to a value of the original thing that should be replaced.
// This function is to be pared with the `With()` function, as in
// `Replace(&original).With(replacement)`.
func Replace[T any](original *T) replacer[T] {
	return replacer[T]{
		original: original,
	}
}

// With can be called on a replacer to set the replacement.
// It returns the interface that can be used in the `Run()` function.
// This function is to be pared with the `Replace()` function, as in
// `Replace(&original).With(replacement)`.
func (f replacer[T]) With(replacement T) Replacer {
	f.replacement = replacement
	return &f
}

// inject caches the original and overwrites it with the replacement.
func (f *replacer[T]) inject() {
	f.origCache = *f.original
	*f.original = f.replacement
}

// revert sets the original back to the value we cached in `inject()`.
func (f *replacer[T]) revert() {
	*f.original = f.origCache
}
