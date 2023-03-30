package CreateUpdateDeleteContainerRegistry

import (
	"automation-suite/PipelineConfigRouter"
	Base "automation-suite/testUtils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strconv"
)

type CreateContainerRegistryFlowsTestSuite struct {
	suite.Suite
	authToken string
}

func (suite *CreateContainerRegistryFlowsTestSuite) SetupSuite() {
	suite.authToken = Base.GetAuthToken()
}

func (suite *CreateContainerRegistryFlowsTestSuite) checkForCiLogs(pipelineId string, ciWorkflowId string, LogLineNumber int, logString string) {
	workflowsDownloadUrl := PipelineConfigRouter.GetCiPipelineBaseUrl + "/" + pipelineId + "/workflows"
	workflows, err := PipelineConfigRouter.FetchCiWorkflows(workflowsDownloadUrl, suite.authToken)
	assert.True(suite.T(), err == nil, err)
	workflowPodName := ""
	workflowNameSpace := ""
	workflowResponses := workflows.Result
	for _, response := range workflowResponses {
		if ciWorkflowId == strconv.Itoa(response.Id) {
			workflowPodName = response.Name
			workflowNameSpace = response.Namespace
		}
	}
	ciLogsDownloadUrlFormat := PipelineConfigRouter.GetCiPipelineBaseUrl + "/%s/workflow/%s/logs"
	ciLogsDownloadUrl := fmt.Sprintf(ciLogsDownloadUrlFormat, pipelineId, ciWorkflowId)
	suite.hitAndCheckBuildLogs(ciLogsDownloadUrl, workflowNameSpace, workflowPodName, LogLineNumber, logString)
}

func (suite *CreateContainerRegistryFlowsTestSuite) hitAndCheckBuildLogs(downloadUrl string, namespace string, wfName string, logLineIndex int, logString string) {
	PipelineConfigRouter.HitLogsDownloadApi(downloadUrl, suite.authToken, suite.T(), logLineIndex, logString)
}
