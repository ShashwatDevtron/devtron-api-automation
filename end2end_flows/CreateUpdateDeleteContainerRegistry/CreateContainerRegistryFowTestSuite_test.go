package CreateUpdateDeleteContainerRegistry

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCreateContainerRegistryFlowsSuite(t *testing.T) {
	suite.Run(t, new(CreateContainerRegistryFlowsTestSuite))
}
