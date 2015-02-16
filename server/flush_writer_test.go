package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type FlushWriterTestSuite struct {
	suite.Suite
}

func (suite *FlushWriterTestSuite) TestWrites() {
	w := httptest.NewRecorder()
	subject := newFlushWriter(w)
	subject.Write([]byte("HTTP/1.1 200 OK"))

	assert.Equal(suite.T(), "HTTP/1.1 200 OK", w.Body.String())
}

func (suite *FlushWriterTestSuite) TestWriteFlushesAutomatically() {
	w := httptest.NewRecorder()
	subject := newFlushWriter(w)
	subject.Write([]byte("test"))

	assert.True(suite.T(), w.Flushed)
}

func TestFlushWriterTestSuite(t *testing.T) {
	suite.Run(t, new(FlushWriterTestSuite))
}
