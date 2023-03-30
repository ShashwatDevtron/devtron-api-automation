package GitProviderRouter

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestGitProRouterSuite(t *testing.T) {
	suite.Run(t, new(GitProRouterTestSuite))
}
