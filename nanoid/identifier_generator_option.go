package nanoid

// IdentifierGeneratorOption represents an option for configure IdentifierGenerator object.
type IdentifierGeneratorOption interface {
	apply(generator *IdentifierGenerator)
}

type identifierGeneratorOptionFunc func(generator *IdentifierGenerator)

func (fn identifierGeneratorOptionFunc) apply(generator *IdentifierGenerator) {
	fn(generator)
}

// DefaultAlphabet is the alphabet used for ID characters by default.
const DefaultAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

// WithAlphabet sets up characters set for identifier.
func WithAlphabet(alphabet string) IdentifierGeneratorOption {
	return identifierGeneratorOptionFunc(func(generator *IdentifierGenerator) {
		generator.alphabet = alphabet
	})
}

// DefaultSize is the default size of identifier.
const DefaultSize = 64

// WithSize sets up identifier length.
func WithSize(size int) IdentifierGeneratorOption {
	return identifierGeneratorOptionFunc(func(generator *IdentifierGenerator) {
		generator.size = size
	})
}
