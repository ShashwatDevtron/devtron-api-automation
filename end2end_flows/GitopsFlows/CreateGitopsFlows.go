package GitopsFlows

import (
	"automation-suite/GitopsConfigRouter"
	"automation-suite/PipelineConfigRouter"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"time"
)

func (suite *GitopsFlowsTestSuite) TestGitopsFlows() {
	suite.Run("Testing-Gitops-Flows", func() {
		envConf := Base.ReadBaseEnvConfig()
		gitopsConfig := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)

		createGitopsConfigRequestDto := GitopsConfigRouter.GetGitopsConfigRequestDto(gitopsConfig.Provider, gitopsConfig.GitUsername, gitopsConfig.Host, gitopsConfig.GitToken, gitopsConfig.GitHubOrgId)
		byteValueOfCreateGitopsConfig, _ := json.Marshal(createGitopsConfigRequestDto)

		ValidatedResponse := ValidateGitops(byteValueOfCreateGitopsConfig, gitopsConfig.Provider, gitopsConfig.GitUsername, gitopsConfig.Host, gitopsConfig.GitToken, gitopsConfig.GitHubOrgId, suite.authToken)
		assert.Equal(suite.T(), 200, ValidatedResponse.Code)
		assert.Equal(suite.T(), "Create Repo", ValidatedResponse.Result.SuccessfulStages[0])

		createLinkResponseDto := SaveGitopsConfig(byteValueOfCreateGitopsConfig, gitopsConfig.Provider, gitopsConfig.GitUsername, gitopsConfig.Host, gitopsConfig.GitToken, gitopsConfig.GitHubOrgId, suite.authToken)
		assert.Equal(suite.T(), 200, createLinkResponseDto.Code)
		assert.Equal(suite.T(), "Create Repo", createLinkResponseDto.Result.SuccessfulStages[0])
		createAppApiResponse, workflowResponse := PipelineConfigRouter.CreateNewAppWithCiCd(suite.authToken)
		appCdPipelineResponse := PipelineConfigRouter.HitGetAppCdPipeline(strconv.Itoa(workflowResponse.Result.AppId), suite.authToken)
		time.Sleep(2 * time.Second)
		log.Println("=== Here we are getting pipeline material ===")
		pipelineMaterial := PipelineConfigRouter.HitGetCiPipelineMaterial(workflowResponse.Result.CiPipelines[0].Id, suite.authToken)

		//here we are hitting GetWorkFlow API 2 time one just after the triggerCiPipeline and one after 4 minutes of triggering

		payloadForTriggerCiPipeline := PipelineConfigRouter.CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.Result.CiPipelines[0].Id, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		triggerCiPipelineResponse := PipelineConfigRouter.HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		if triggerCiPipelineResponse.Result.AuthStatus != "allowed for all pipelines" {
			time.Sleep(5 * time.Second)
			triggerCiPipelineResponse = PipelineConfigRouter.HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
			assert.Equal(suite.T(), "allowed for all pipelines", triggerCiPipelineResponse.Result.AuthStatus)
			assert.NotNil(suite.T(), triggerCiPipelineResponse.Result.ApiResponse)
		}
		time.Sleep(5 * time.Second)
		log.Println("=== Here we are getting workflow after triggering ===")
		workflowStatus := PipelineConfigRouter.HitGetWorkflowStatus(createAppApiResponse.Id, suite.authToken)
		if workflowStatus.Result.CiWorkflowStatus[0].CiStatus == "Starting" {
			time.Sleep(5 * time.Second)
			workflowStatus = PipelineConfigRouter.HitGetWorkflowStatus(createAppApiResponse.Id, suite.authToken)
			assert.Equal(suite.T(), "Running", workflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		} else {
			assert.Equal(suite.T(), "Running", workflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		}
		log.Println("=== Here we are getting workflow and verifying the status after triggering via poll function ===")
		assert.True(suite.T(), PipelineConfigRouter.PollForGettingCdDeployStatusAfterTrigger(createAppApiResponse.Id, suite.authToken))
		updatedWorkflowStatus := PipelineConfigRouter.HitGetWorkflowStatus(createAppApiResponse.Id, suite.authToken)
		assert.Equal(suite.T(), "Succeeded", updatedWorkflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		assert.Equal(suite.T(), "Succeeded", updatedWorkflowStatus.Result.CdWorkflowStatus[0].DeployStatus)
		queryparams := map[string]string{"app-id": strconv.Itoa(workflowResponse.Result.AppId), "env-id": strconv.Itoa(appCdPipelineResponse.Result.Pipelines[0].EnvironmentId)}
		appStatus := PipelineConfigRouter.GetAppDetail(queryparams, suite.authToken)
		assert.Equal(suite.T(), "Healthy", appStatus.Result.ResourceTree.Status)

	})
}
