package K8sPodsRouter

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestK8sPodsRouterSuite(t *testing.T) {
	suite.Run(t, new(K8sPodsTestSuite))
}
