package PipelineConfigRouter

import (
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"time"
)

func (suite *PipelinesConfigRouterTestSuite) TestClassD3TriggerCiPipeline() {
	createAppApiResponse, workflowResponse := CreateNewAppWithCiCd(suite.authToken)
	time.Sleep(2 * time.Second)
	log.Println("=== Here we are getting pipeline material ===")
	pipelineMaterial := HitGetCiPipelineMaterial(workflowResponse.Result.CiPipelines[0].Id, suite.authToken)

	//here we are hitting GetWorkFlow API 2 time one just after the triggerCiPipeline and one after 4 minutes of triggering
	suite.Run("A=1=TriggerCiPipelineWithValidPayload", func() {
		payloadForTriggerCiPipeline := CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.Result.CiPipelines[0].Id, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		triggerCiPipelineResponse := HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		if triggerCiPipelineResponse.Result.AuthStatus != "allowed for all pipelines" {
			time.Sleep(5 * time.Second)
			triggerCiPipelineResponse = HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
			assert.Equal(suite.T(), "allowed for all pipelines", triggerCiPipelineResponse.Result.AuthStatus)
			assert.NotNil(suite.T(), triggerCiPipelineResponse.Result.ApiResponse)
		}
		time.Sleep(5 * time.Second)
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
		assert.True(suite.T(), PollForGettingCdDeployStatusAfterTrigger(createAppApiResponse.Id, suite.authToken))
		updatedWorkflowStatus := HitGetWorkflowStatus(createAppApiResponse.Id, suite.authToken)
		assert.Equal(suite.T(), "Succeeded", updatedWorkflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		assert.Equal(suite.T(), "Succeeded", updatedWorkflowStatus.Result.CdWorkflowStatus[0].DeployStatus)
	})

	suite.Run("A=2=TriggerCiPipelineWithInvalidateCacheAsFalse", func() {
		payloadForTriggerCiPipeline := CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.Result.CiPipelines[0].Id, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		triggerCiPipelineResponse := HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		if triggerCiPipelineResponse.Result.AuthStatus != "allowed for all pipelines" {
			time.Sleep(2 * time.Second)
			triggerCiPipelineResponse = HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
			time.Sleep(5 * time.Second)
			assert.Equal(suite.T(), "allowed for all pipelines", triggerCiPipelineResponse.Result.AuthStatus)
			assert.NotNil(suite.T(), triggerCiPipelineResponse.Result.ApiResponse)
		}
	})

	suite.Run("A=3=TriggerCiPipelineWithInvalidPipelineId", func() {
		invalidPipeLineId := Base.GetRandomNumberOf9Digit()
		payloadForTriggerCiPipeline := CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, invalidPipeLineId, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		triggerCiPipelineResponse := HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		assert.Equal(suite.T(), "pg: no rows in result set", triggerCiPipelineResponse.Errors[0].UserMessage)
	})

	suite.Run("A=4=TriggerCiPipelineWithInvalidPipelineId", func() {
		invalidMaterialId := Base.GetRandomNumberOf9Digit()
		payloadForTriggerCiPipeline := CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.Result.CiPipelines[0].Id, invalidMaterialId, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		triggerCiPipelineResponse := HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		assert.Equal(suite.T(), "[{pg: no rows in result set}]", triggerCiPipelineResponse.Errors[0].InternalMessage)
	})

	DeleteAppWithCiCd(suite.authToken)
}

func PollForGettingCdDeployStatusAfterTrigger(id int, authToken string) bool {
	count := 0
	for {
		updatedWorkflowStatus := HitGetWorkflowStatus(id, authToken)
		deploymentStatus := updatedWorkflowStatus.Result.CdWorkflowStatus[0].DeployStatus
		time.Sleep(1 * time.Second)
		count = count + 1
		if deploymentStatus == "Succeeded" || count >= 600 {
			break
		}
	}
	return true
}
func PollForGettingCdAppStatusAfterTrigger(id int, authToken string, suite *PipelinesConfigRouterTestSuite) bool {
	count := 0
	for {
		apiResponse := GetAppDeploymentStatusTimeline(id, 1, authToken)
		time.Sleep(1 * time.Second)
		count = count + 1
		if apiResponse.Code == 404 {
			log.Println("AppTimelineStatus====> Waiting For Deployment")
		}
		//write the test cases here
		if len(apiResponse.Result.Timelines) != 0 {
			if len(apiResponse.Result.Timelines) >= 1 {
				suite.Run("TestDeploymentInitiation", func() {
					assert.NotEqual(suite.T(), nil, apiResponse)
					assert.Equal(suite.T(), 200, apiResponse.Code)
					assert.Equal(suite.T(), 0, len(apiResponse.Error))
					assert.NotEqual(suite.T(), nil, apiResponse.Result)
					isDeploymentStarted := len(apiResponse.Result.Timelines) > 0
					assert.Equal(suite.T(), true, isDeploymentStarted)
					assert.Equal(suite.T(), TIMELINE_STATUS_DEPLOYMENT_INITIATED, apiResponse.Result.Timelines[0].Status)
				})
			}
			if len(apiResponse.Result.Timelines) >= 2 {
				suite.Run("TestDeploymentGitCommitInitiation", func() {
					assert.NotEqual(suite.T(), nil, apiResponse)
					assert.Equal(suite.T(), 200, apiResponse.Code)
					assert.Equal(suite.T(), 0, len(apiResponse.Error))
					assert.NotEqual(suite.T(), nil, apiResponse.Result)
					isDeploymentStarted := len(apiResponse.Result.Timelines) > 1
					assert.Equal(suite.T(), true, isDeploymentStarted)
					assert.Equal(suite.T(), TIMELINE_STATUS_DEPLOYMENT_INITIATED, apiResponse.Result.Timelines[0].Status)
					gitStatus := apiResponse.Result.Timelines[1].Status == TIMELINE_STATUS_GIT_COMMIT || apiResponse.Result.Timelines[1].Status == TIMELINE_STATUS_GIT_COMMIT_FAILED
					assert.Equal(suite.T(), true, gitStatus)
				})
			}
			if len(apiResponse.Result.Timelines) >= 3 {
				suite.Run("TestGitCommitSuccessAndKubectlStarted", func() {
					assert.NotEqual(suite.T(), nil, apiResponse)
					assert.Equal(suite.T(), 200, apiResponse.Code)
					assert.Equal(suite.T(), 0, len(apiResponse.Error))
					assert.NotEqual(suite.T(), nil, apiResponse.Result)
					isDeploymentStarted := len(apiResponse.Result.Timelines) > 2
					assert.Equal(suite.T(), true, isDeploymentStarted)
					assert.Equal(suite.T(), TIMELINE_STATUS_DEPLOYMENT_INITIATED, apiResponse.Result.Timelines[0].Status)
					assert.Equal(suite.T(), TIMELINE_STATUS_GIT_COMMIT, apiResponse.Result.Timelines[1].Status)
					kubectlStatus := apiResponse.Result.Timelines[2].Status == TIMELINE_STATUS_KUBECTL_APPLY_STARTED
					assert.Equal(suite.T(), true, kubectlStatus)
				})
			}
			if len(apiResponse.Result.Timelines) >= 4 {
				suite.Run("TestGitCommitSuccessAndKubectlApply", func() {
					assert.NotEqual(suite.T(), nil, apiResponse)
					assert.Equal(suite.T(), 200, apiResponse.Code)
					assert.Equal(suite.T(), 0, len(apiResponse.Error))
					assert.NotEqual(suite.T(), nil, apiResponse.Result)
					isDeploymentStarted := len(apiResponse.Result.Timelines) > 3
					assert.Equal(suite.T(), true, isDeploymentStarted)
					assert.Equal(suite.T(), TIMELINE_STATUS_DEPLOYMENT_INITIATED, apiResponse.Result.Timelines[0].Status)
					assert.Equal(suite.T(), TIMELINE_STATUS_GIT_COMMIT, apiResponse.Result.Timelines[1].Status)
					assert.Equal(suite.T(), TIMELINE_STATUS_KUBECTL_APPLY_STARTED, apiResponse.Result.Timelines[2].Status)
					kubectlStatus := apiResponse.Result.Timelines[3].Status == TIMELINE_STATUS_KUBECTL_APPLY_SYNCED
					assert.Equal(suite.T(), true, kubectlStatus)
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
				})
			}
			if len(apiResponse.Result.Timelines) >= 5 && (apiResponse.Result.Timelines[4].Status == TIMELINE_STATUS_APP_HEALTHY || apiResponse.Result.Timelines[4].Status == TIMELINE_STATUS_FETCH_TIMED_OUT) {
				break
			}
		}
		if count >= 600 {
			break
		}
	}
	return true
}
