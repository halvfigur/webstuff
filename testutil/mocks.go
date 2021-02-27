package testutil

import (
	"github.com/halvfigur/webstuff/model"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

const MockStoreUser = "User"

func (m *MockStore) User(username string) (model.User, error) {
	args := m.Called(username)
	return args.Get(0).(model.User), args.Error(1)
}

const MockSessionCodecEncode = "Encode"

type MockSessionCodec struct {
	mock.Mock
}

func (m *MockSessionCodec) Encode(session []byte) (string, error) {
	args := m.Called(session)
	return args.String(0), args.Error(1)
}

const MockSessionCodecDecode = "Decode"

func (m *MockSessionCodec) Decode(token string) ([]byte, error) {
	args := m.Called(token)
	return args.Get(0).([]byte), args.Error(1)
}
