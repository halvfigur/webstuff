package controller

import (
	"crypto/rand"
	"encoding/json"
	"testing"
	"time"

	"github.com/halvfigur/webstuff/model"
	"github.com/stretchr/testify/require"
)

func TestSessionCodec(t *testing.T) {
	expires, _ := time.Parse(time.RFC3339, time.RFC3339)

	session := model.Session{
		Username: "test_user",
		Expires:  expires,
	}

	blob, _ := json.Marshal(session)

	key := make([]byte, 32)
	rand.Read(key)
	codec := NewPasetoSessionCodec(key)

	encoded, err := codec.Encode(blob)
	require.Nil(t, err)

	decoded, err := codec.Decode(encoded)
	require.Nil(t, err)

	var decodedSession model.Session
	err = json.Unmarshal(decoded, &decodedSession)
	require.Nil(t, err)

	require.Equal(t, decodedSession.Username, session.Username)
	require.Equal(t, decodedSession.Expires, session.Expires)

}
