package nanoid

import (
	"context"

	gonanoid "github.com/matoous/go-nanoid/v2"
	banking "github.com/morozovcookie/agat-banking"
	"github.com/pkg/errors"
)

var _ banking.IdentifierGenerator = (*IdentifierGenerator)(nil)

// IdentifierGenerator represents a service for generating unique identifier.
type IdentifierGenerator struct {
	alphabet string
	size     int
}

// NewIdentifierGenerator returns a new IdentifierGenerator instance.
func NewIdentifierGenerator(opts ...IdentifierGeneratorOption) *IdentifierGenerator {
	gen := &IdentifierGenerator{
		alphabet: DefaultAlphabet,
		size:     DefaultSize,
	}

	for _, opt := range opts {
		opt.apply(gen)
	}

	return gen
}

// GenerateIdentifier returns unique identifier.
func (gen *IdentifierGenerator) GenerateIdentifier(_ context.Context) (banking.ID, error) {
	id, err := gonanoid.Generate(gen.alphabet, gen.size)
	if err != nil {
		return banking.EmptyID, errors.Wrap(err, "generate identifier")
	}

	return banking.ID(id), nil
}
