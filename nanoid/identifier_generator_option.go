package nanoid

// IdentifierGeneratorOption represents an option for configure IdentifierGenerator object.
type IdentifierGeneratorOption interface {
	apply(generator *IdentifierGenerator)
}

type identifierGeneratorOptionFunc func(generator *IdentifierGenerator)

func (fn identifierGeneratorOptionFunc) apply(generator *IdentifierGenerator) {
	fn(generator)
}

// WithAlphabet sets up characters set for identifier.
func WithAlphabet(alphabet string) IdentifierGeneratorOption {
	return identifierGeneratorOptionFunc(func(generator *IdentifierGenerator) {
		generator.alphabet = alphabet
	})
}

// WithSize sets up identifier length.
func WithSize(size int) IdentifierGeneratorOption {
	return identifierGeneratorOptionFunc(func(generator *IdentifierGenerator) {
		generator.size = size
	})
}
