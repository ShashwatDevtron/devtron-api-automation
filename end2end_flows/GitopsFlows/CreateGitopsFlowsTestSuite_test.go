package GitopsFlows

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCreateGitopsFlowsSuite(t *testing.T) {
	suite.Run(t, new(GitopsFlowsTestSuite))
}
