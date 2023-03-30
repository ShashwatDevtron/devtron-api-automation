package PipelineConfigRouter

import (
	PipelineConfigRouterResponseDTOs "automation-suite/PipelineConfigRouter/ResponseDTOs"
	Base "automation-suite/testUtils"
	"github.com/stretchr/testify/assert"

	"encoding/json"
	"log"
	"time"
)

var ciTriggerWorkflowIdPtr *string

func (suite *PipelinesConfigRouterTestSuite) TestGetAppDeploymentStatusTimeline() {
	createAppApiResponse, workflowResponse := CreateNewAppWithCiCd(suite.authToken)

	time.Sleep(2 * time.Second)
	log.Println("=== Here we are getting workflow status material ===")
	updatedWorkflowStatus := HitGetWorkflowStatus(createAppApiResponse.Id, suite.authToken)
	if updatedWorkflowStatus.Result.CdWorkflowStatus[0].DeployStatus == "Not Deployed" || updatedWorkflowStatus.Code != 200 {
		log.Println("=== Here we are getting pipeline material ===")
		pipelineMaterial := HitGetCiPipelineMaterial(workflowResponse.Result.CiPipelines[0].Id, suite.authToken)
		log.Println("=== Here we are Triggering CI/CD and verifying CI/CD Deploy Status ===")
		time.Sleep(10 * time.Second)
		TriggerAndVerifyCiPipeline(createAppApiResponse, pipelineMaterial, workflowResponse.Result.CiPipelines[0].Id, suite)
	}
	//write tests with invalid git-ops configuration
	//end of test cases
	//DeleteAppWithCiCd(suite.authToken)
}

func TriggerAndVerifyCiPipeline(createAppApiResponse Base.CreateAppRequestDto, pipelineMaterial PipelineConfigRouterResponseDTOs.GetCiPipelineMaterialResponseDTO, CiPipelineID int, suite *PipelinesConfigRouterTestSuite) string {
	if ciTriggerWorkflowIdPtr != nil {
		return *ciTriggerWorkflowIdPtr
	}

	ciTriggerWorkflowId := ""
	payloadForTriggerCiPipeline := CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, CiPipelineID, pipelineMaterial.Result[0].Id, true)
	bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
	triggerCiPipelineResponse := HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
	if triggerCiPipelineResponse.Result.AuthStatus != "allowed for all pipelines" {
		time.Sleep(2 * time.Second)
		triggerCiPipelineResponse = HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		assert.Equal(suite.T(), "allowed for all pipelines", triggerCiPipelineResponse.Result.AuthStatus)
		assert.NotNil(suite.T(), triggerCiPipelineResponse.Result.ApiResponse)
	}
	ciTriggerWorkflowId = triggerCiPipelineResponse.Result.ApiResponse
	time.Sleep(10 * time.Second)
	log.Println("=== Here we are getting workflow after triggering ===")
	workflowStatus := HitGetWorkflowStatus(createAppApiResponse.Id, suite.authToken)
	if workflowStatus.Result.CiWorkflowStatus[0].CiStatus == "Starting" {
		time.Sleep(5 * time.Second)
		workflowStatus = HitGetWorkflowStatus(createAppApiResponse.Id, suite.authToken)
		assert.Equal(suite.T(), "Running", workflowStatus.Result.CiWorkflowStatus[0].CiStatus)
	} else {
		assert.Equal(suite.T(), "Running", workflowStatus.Result.CiWorkflowStatus[0].CiStatus)
	}
	log.Println("=== Here we are getting workflow and verifying the status after triggering via poll function ===")
	assert.True(suite.T(), PollForGettingCdAppStatusAfterTrigger(createAppApiResponse.Id, suite.authToken, suite))
	apiResponse := GetAppDeploymentStatusTimeline(createAppApiResponse.Id, 1, suite.authToken)
	time.Sleep(2 * time.Second)
	suite.Run("TestHealthyDeploymentStatus", func() {
		assert.NotEqual(suite.T(), nil, apiResponse)
		assert.Equal(suite.T(), 200, apiResponse.Code)
		assert.Equal(suite.T(), 0, len(apiResponse.Error))
		assert.NotEqual(suite.T(), nil, apiResponse.Result)
		isDeploymentStarted := len(apiResponse.Result.Timelines) == 5
		assert.Equal(suite.T(), true, isDeploymentStarted)
		assert.Equal(suite.T(), TIMELINE_STATUS_DEPLOYMENT_INITIATED, apiResponse.Result.Timelines[0].Status)
		assert.Equal(suite.T(), TIMELINE_STATUS_GIT_COMMIT, apiResponse.Result.Timelines[1].Status)
		assert.Equal(suite.T(), TIMELINE_STATUS_KUBECTL_APPLY_STARTED, apiResponse.Result.Timelines[2].Status)
		assert.Equal(suite.T(), TIMELINE_STATUS_KUBECTL_APPLY_SYNCED, apiResponse.Result.Timelines[3].Status)
		isAtleastOneK8sObjectPresent := len(apiResponse.Result.Timelines[2].ResourceDetails) > 0
		assert.Equal(suite.T(), true, isAtleastOneK8sObjectPresent)
		for _, resource := range apiResponse.Result.Timelines[2].ResourceDetails {
			assert.NotNil(suite.T(), resource.Id)
			assert.NotNil(suite.T(), resource.ResourceStatus)
			assert.NotNil(suite.T(), resource.ResourceKind)
			assert.NotNil(suite.T(), resource.ResourceGroup)
			assert.NotNil(suite.T(), resource.ResourceName)
			assert.NotNil(suite.T(), resource.ResourcePhase)
		}
		assert.Equal(suite.T(), TIMELINE_STATUS_APP_HEALTHY, apiResponse.Result.Timelines[4].Status)
	})
	ciTriggerWorkflowIdPtr = &ciTriggerWorkflowId
	return ciTriggerWorkflowId
}
