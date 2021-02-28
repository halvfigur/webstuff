package webstuff

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSessionCodec(t *testing.T) {
	session := []byte("secret session stuff")

	key := make([]byte, 32)
	rand.Read(key)
	codec := NewPasetoSessionCodec(key)

	encoded, err := codec.Encode(session)
	require.Nil(t, err)

	decodedSession, err := codec.Decode(encoded)
	require.Nil(t, err)

	require.Equal(t, session, decodedSession)
}
