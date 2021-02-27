package controller

import (
	"github.com/o1egl/paseto"
)

// SessionCodec encodes and decodes session cookies
type SessionCodec interface {
	// Encode a session cookie
	Encode(session []byte) (string, error)
	// Decode a session cookie
	Decode(token string) ([]byte, error)
}

// PasetoSessionCodec provides PASETO security tokens
type PasetoSessionCodec struct {
	v2  *paseto.V2
	key []byte
}

// NewPasetoSessionCodec creates a new PasetoSessionCodec using key
func NewPasetoSessionCodec(key []byte) *PasetoSessionCodec {
	return &PasetoSessionCodec{
		v2:  paseto.NewV2(),
		key: key,
	}
}

// Encode a session cookie
func (c *PasetoSessionCodec) Encode(session []byte) (string, error) {
	return c.v2.Encrypt(c.key, session, nil)
}

// Decode a session cookie
func (c *PasetoSessionCodec) Decode(token string) ([]byte, error) {
	var session []byte

	if err := c.v2.Decrypt(token, c.key, &session, nil); err != nil {
		return nil, err
	}

	return session, nil
}
