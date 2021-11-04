package json

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/morozovcookie/agat-banking/mock"
	"github.com/stretchr/testify/assert"
)

func TestSecretString_DecryptedString(t *testing.T) {
	type meta struct {
		name    string
		enabled bool
	}
	type fields struct {
		setupSecretStringFn func() *mock.SecretString
	}
	type wants struct {
		decryptedString string
	}

	tests := []struct {
		meta   meta
		fields fields
		wants  wants
	}{
		{
			meta: meta{
				name:    "pass",
				enabled: true,
			},
			fields: fields{
				setupSecretStringFn: func() *mock.SecretString {
					ss := mock.NewSecretString()

					ss.On("DecryptedString").
						Return("decrypted string")

					return ss
				},
			},
			wants: wants{
				decryptedString: "decrypted string",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.meta.name, func(t *testing.T) {
			if !tt.meta.enabled {
				t.SkipNow()
			}

			ss := tt.fields.setupSecretStringFn()

			decryptedString := NewSecretString(ss, nil).DecryptedString()

			assert.Equal(t, tt.wants.decryptedString, decryptedString)

			ss.AssertExpectations(t)
		})
	}
}

func TestSecretString_EncryptedString(t *testing.T) {
	type meta struct {
		name    string
		enabled bool
	}
	type fields struct {
		setupSecretStringFn func() *mock.SecretString
	}
	type wants struct {
		encryptedString string
	}

	tests := []struct {
		meta   meta
		fields fields
		wants  wants
	}{
		{
			meta: meta{
				name:    "pass",
				enabled: true,
			},
			fields: fields{
				setupSecretStringFn: func() *mock.SecretString {
					ss := mock.NewSecretString()

					ss.On("EncryptedString").
						Return("encrypted string")

					return ss
				},
			},
			wants: wants{
				encryptedString: "encrypted string",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.meta.name, func(t *testing.T) {
			if !tt.meta.enabled {
				t.SkipNow()
			}

			ss := tt.fields.setupSecretStringFn()

			encryptedString := NewSecretString(ss, nil).EncryptedString()

			assert.Equal(t, tt.wants.encryptedString, encryptedString)

			ss.AssertExpectations(t)
		})
	}
}

func TestSecretString_String(t *testing.T) {
	type meta struct {
		name    string
		enabled bool
	}
	type fields struct {
		setupSecretStringFn func() *mock.SecretString
	}
	type wants struct {
		str string
	}

	tests := []struct {
		meta   meta
		fields fields
		wants  wants
	}{
		{
			meta: meta{
				name:    "pass",
				enabled: true,
			},
			fields: fields{
				setupSecretStringFn: func() *mock.SecretString {
					ss := mock.NewSecretString()

					ss.On("String").
						Return("just a string")

					return ss
				},
			},
			wants: wants{
				str: "just a string",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.meta.name, func(t *testing.T) {
			if !tt.meta.enabled {
				t.SkipNow()
			}

			ss := tt.fields.setupSecretStringFn()

			str := NewSecretString(ss, nil).String()

			assert.Equal(t, tt.wants.str, str)

			ss.AssertExpectations(t)
		})
	}
}

func TestSecretString_MarshalJSON(t *testing.T) {
	type meta struct {
		name    string
		enabled bool
	}
	type fields struct {
		setupSecretStringFn func() *mock.SecretString
	}
	type wants struct {
		bytes []byte
		err   bool
	}

	tests := []struct {
		meta   meta
		fields fields
		wants  wants
	}{
		{
			meta: meta{
				name:    "pass",
				enabled: true,
			},
			fields: fields{
				setupSecretStringFn: func() *mock.SecretString {
					ss := mock.NewSecretString()

					ss.On("DecryptedString").
						Return("decrypted string")

					return ss
				},
			},
			wants: wants{
				bytes: bytes.NewBufferString(`"decrypted string"` + "\n").Bytes(),
				err:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.meta.name, func(t *testing.T) {
			if !tt.meta.enabled {
				t.SkipNow()
			}

			var (
				ss           = tt.fields.setupSecretStringFn()
				secretString = NewSecretString(ss, nil)

				buf = new(bytes.Buffer)
				err = json.NewEncoder(buf).Encode(secretString)
			)

			assert.Equal(t, tt.wants.bytes, buf.Bytes())
			assert.Equal(t, tt.wants.err, err != nil)

			ss.AssertExpectations(t)
		})
	}
}

func TestSecretString_UnmarshalJSON(t *testing.T) {
	type meta struct {
		name    string
		enabled bool
	}
	type fields struct {
		setupSecretFactoryFn func() *mock.SecretFactory
	}
	type args struct {
		bytes []byte
	}
	type wants struct {
		setupSecretStringFn func() *mock.SecretString
		err                 bool
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
				setupSecretFactoryFn: func() *mock.SecretFactory {
					var (
						factory = mock.NewSecretFactory()
						ss      = mock.NewSecretString()
					)

					ss.On("String").
						Return("just a string")
					ss.On("DecryptedString").
						Return("decrypted string")
					ss.On("EncryptedString").
						Return("encrypted string")

					factory.On("CreateFromDecryptedData", bytes.NewBufferString("decrypted string")).
						Return(ss, (error)(nil))

					return factory
				},
			},
			args: args{
				bytes: bytes.NewBufferString(`"decrypted string" + "\n"`).Bytes(),
			},
			wants: wants{
				setupSecretStringFn: func() *mock.SecretString {
					ss := mock.NewSecretString()

					ss.On("String").
						Return("just a string")
					ss.On("DecryptedString").
						Return("decrypted string")
					ss.On("EncryptedString").
						Return("encrypted string")

					return ss
				},
				err: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.meta.name, func(t *testing.T) {
			if !tt.meta.enabled {
				t.SkipNow()
			}

			var (
				factory      = tt.fields.setupSecretFactoryFn()
				secretString = NewSecretString(nil, factory)

				err = json.NewDecoder(bytes.NewBuffer(tt.args.bytes)).Decode(secretString)
			)

			assert.Equal(t, tt.wants.err, err != nil)
			factory.AssertExpectations(t)

			if err != nil {
				return
			}

			ss := tt.wants.setupSecretStringFn()

			assert.Equal(t, ss.String(), secretString.String())
			assert.Equal(t, ss.DecryptedString(), secretString.DecryptedString())
			assert.Equal(t, ss.EncryptedString(), secretString.EncryptedString())
		})
	}
}
