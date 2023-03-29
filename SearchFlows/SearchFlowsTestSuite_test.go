package SearchFlows

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestSearchFlowsSuite(t *testing.T) {
	suite.Run(t, new(SearchFlowTestSuite))
}
