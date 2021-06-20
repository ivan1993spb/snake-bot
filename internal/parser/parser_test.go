package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
	Parser *Parser
}

func (suite *ParserTestSuite) SetupTest() {
	suite.Parser = &Parser{}
}

var messageBroadcastPayload = []byte(`"user joined your game group"`)

var messageBroadcastPayloadCorrupted = []byte(`user joined your game group"`)

func (suite *ParserTestSuite) TestParseBroadcast() {
	printer := new(MockPrinter)
	suite.Parser.Printer = printer
	printer.On("Print", "broadcast", "message", "user joined your game group").Return()
	err := suite.Parser.ParseBroadcast(messageBroadcastPayload)
	assert.Nil(suite.T(), err)
	printer.AssertExpectations(suite.T())
}

func (suite *ParserTestSuite) TestParseBroadcastCorruptedPayload() {
	err := suite.Parser.ParseBroadcast(messageBroadcastPayloadCorrupted)
	assert.NotNil(suite.T(), err)
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
