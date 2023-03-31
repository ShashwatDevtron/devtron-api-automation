package K8sPodsRouter

import (
	Base "automation-suite/testUtils"
	"automation-suite/testdata/testUtils"
	"github.com/stretchr/testify/suite"
)

type K8sPodsTestSuite struct {
	suite.Suite
	authToken string
}

func (suite *K8sPodsTestSuite) SetupSuite() {
	suite.authToken = Base.GetAuthToken()
}

func (suite *K8sPodsTestSuite) HitCheckLogsApi(name string, containerName string, containerId string, namespace string, acdAppId string, envId string, follow string, tailLines string) {
	ciLogsDownloadUrl := K8sPodsRouterBaseUrl + name + "?containerName=" + containerName + "&clusterId=" + containerId + "&namespace=" + namespace + "&acdAppId=" + acdAppId + "&envId=" + envId + "&follow=" + follow + "&tailLines=" + tailLines
	testUtils.ReadEventStreamsForSpecificApiAndVerifyResult(ciLogsDownloadUrl, suite.authToken, suite.T(), 0, "Listening on", true)
}
