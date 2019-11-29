package buffer

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

func (m *Mock) Write(item interface{}) {
	m.Called(item)
}
