package parser

import (
	"github.com/stretchr/testify/mock"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

type MockCountdown struct {
	mock.Mock
}

func (m *MockCountdown) Countdown(sec int) {
	m.Called(sec)
}

type MockMe struct {
	mock.Mock
}

func (m *MockMe) Me(id uint32) {
	m.Called(id)
}

type MockSize struct {
	mock.Mock
}

func (m *MockSize) Size(width, height uint8) {
	m.Called(width, height)
}

type MockGame struct {
	mock.Mock
}

func (m *MockGame) Create(object *types.Object) {
	m.Called(object)
}

func (m *MockGame) Update(object *types.Object) {
	m.Called(object)
}

func (m *MockGame) Delete(object *types.Object) {
	m.Called(object)
}

type MockPrinter struct {
	mock.Mock
}

func (m *MockPrinter) Print(level, what, message string) {
	m.Called(level, what, message)
}
