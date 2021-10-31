package aes_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/morozovcookie/agat-banking/aes"
	"github.com/morozovcookie/agat-banking/mock"
	"github.com/stretchr/testify/assert"
)

func TestSecretFactory_CreateFromDecryptedData(t *testing.T) {
	type meta struct {
		name    string
		enabled bool
	}
	type fields struct {
		key string

		setupNonceGenerator func() *mock.NonceGenerator
	}
	type args struct {
		ctx           context.Context
		decryptedText string
	}
	type wants struct {
		encryptedString string
		decryptedString string

		err bool
	}

	tests := []struct {
		meta   meta
		fields fields
		args   args
		wants  wants
	}{
		{
			meta: meta{
				name:    "pass",
				enabled: true,
			},
			fields: fields{
				key: "a7681ff138d941377c55aefb4ab667b833a823e582c91317f5b5e33c09e6891e",

				setupNonceGenerator: func() *mock.NonceGenerator {
					generator := &mock.NonceGenerator{}

					generator.On("GenerateNonce", 12).
						Return([]byte{'a', 'g', 'g', '3', 'k', 's', 'v', 't', 'l', 'a', 'm', '7'}, (error)(nil))

					return generator
				},
			},
			args: args{
				ctx:           context.Background(),
				decryptedText: "super-secret-string",
			},
			wants: wants{
				encryptedString: "616767336b7376746c616d376276c58836970ec0fd86c63d74233cfcea0f36a3dd83247363b1b850152" +
					"b867e1b51c7",
				decryptedString: "super-secret-string",

				err: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.meta.name, func(t *testing.T) {
			if !tt.meta.enabled {
				t.SkipNow()
			}

			nonceGenerator := tt.fields.setupNonceGenerator()

			factory, err := aes.NewSecretFactory(nonceGenerator, bytes.NewBufferString(tt.fields.key))
			assert.NoError(t, err)

			ss, err := factory.CreateFromDecryptedData(tt.args.ctx, bytes.NewBufferString(tt.args.decryptedText))
			assert.Equal(t, tt.wants.encryptedString, ss.EncryptedString())
			assert.Equal(t, tt.wants.decryptedString, ss.DecryptedString())
			assert.Equal(t, tt.wants.err, err != nil)

			nonceGenerator.AssertExpectations(t)
		})
	}
}

func TestSecretFactory_CreateFromEncryptedData(t *testing.T) {
	type meta struct {
		name    string
		enabled bool
	}
	type fields struct {
		key string

		setupNonceGenerator func() *mock.NonceGenerator
	}
	type args struct {
		ctx           context.Context
		encryptedText string
	}
	type wants struct {
		encryptedString string
		decryptedString string

		err bool
	}

	tests := []struct {
		meta   meta
		fields fields
		args   args
		wants  wants
	}{
		{
			meta: meta{
				name:    "pass",
				enabled: true,
			},
			fields: fields{
				key: "a7681ff138d941377c55aefb4ab667b833a823e582c91317f5b5e33c09e6891e",
			},
			args: args{
				ctx: context.Background(),
				encryptedText: "616767336b7376746c616d376276c58836970ec0fd86c63d74233cfcea0f36a3dd83247363b1b850152b8" +
					"67e1b51c7",
			},
			wants: wants{
				encryptedString: "616767336b7376746c616d376276c58836970ec0fd86c63d74233cfcea0f36a3dd83247363b1b850152" +
					"b867e1b51c7",
				decryptedString: "super-secret-string",

				err: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.meta.name, func(t *testing.T) {
			if !tt.meta.enabled {
				t.SkipNow()
			}

			factory, err := aes.NewSecretFactory(nil, bytes.NewBufferString(tt.fields.key))
			assert.NoError(t, err)

			ss, err := factory.CreateFromEncryptedData(tt.args.ctx, bytes.NewBufferString(tt.args.encryptedText))
			assert.Equal(t, tt.wants.encryptedString, ss.EncryptedString())
			assert.Equal(t, tt.wants.decryptedString, ss.DecryptedString())
			assert.Equal(t, tt.wants.err, err != nil)
		})
	}
}
