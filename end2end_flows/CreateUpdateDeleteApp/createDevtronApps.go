package CreateUpdateDeleteApp

import (
	"automation-suite/HelperRouter"
	"automation-suite/PipelineConfigRouter"
	dtos "automation-suite/PipelineConfigRouter/ResponseDTOs"
	"automation-suite/RbacFlows"
	"automation-suite/RbacFlows/RequestDTOs"
	"automation-suite/TeamRouter"
	TeamRouterResponseDTOs "automation-suite/TeamRouter/ResponseDTOs"
	"automation-suite/testUtils"
	//"automation-suite/testdata/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/sjson"
	"log"
	"strconv"
	"strings"
	"time"
)

func (suite *CreateDevtronAppFlowsTestSuite) TestFlowsForDevtronApps() {
	var (
		createAppApiResponsePtr *PipelineConfigRouter.CreateAppRequestDto
		workflowResponsePtr     *dtos.CreateWorkflowResponseDto
		savePipelineResponsePtr *dtos.SaveCdPipelineResponseDTO
	)
	suite.Run("A0=CreateDevtronAppWithoutTags", func() {
		var devtronDeletion RequestDTOs.RbacDevtronDeletion
		fetchAllTeamResponseDto := GetAllTeamsAutocomplete(suite.authToken)
		log.Println("Validating the response of FetchAllTeam API")
		var responseOfCreateProject TeamRouterResponseDTOs.SaveTeamResponseDTO
		if len(fetchAllTeamResponseDto.Result) == 0 {
			saveTeamRequestDto := TeamRouter.GetSaveTeamRequestDto()
			byteValueOfStruct, _ := json.Marshal(saveTeamRequestDto)
			responseOfCreateProject = RbacFlows.CreateProject(byteValueOfStruct, suite.authToken)
			devtronDeletion.ProjectPayload, _ = json.Marshal(responseOfCreateProject.Result)
			assert.Equal(suite.T(), 200, responseOfCreateProject.Code)
			assert.Equal(suite.T(), saveTeamRequestDto.Name, responseOfCreateProject.Result.Name)
		}
		queryParams := map[string]string{"auth": "true"}
		allEnvironmentDetailsResponse := GetAllEnvsAutocomplete(queryParams, suite.authToken)
		log.Println("Validating the response of GetAllEnvironmentDetails API")
		assert.NotNil(suite.T(), allEnvironmentDetailsResponse.Result)
		assert.Equal(suite.T(), 200, allEnvironmentDetailsResponse.Code)
		if len(allEnvironmentDetailsResponse.Result) == 0 {
			saveEnvRequestDto := RbacFlows.GetSaveEnvRequestDto()
			saveEnvRequestDto.EnvironmentIdentifier = "default_cluster__" + saveEnvRequestDto.Namespace
			byteValueOfStruct, _ := json.Marshal(saveEnvRequestDto)
			responseOfCreateEnvironment := RbacFlows.CreateEnv(byteValueOfStruct, suite.authToken)
			devtronDeletion.EnvPayLoad, _ = json.Marshal(responseOfCreateEnvironment.Result)
			assert.Equal(suite.T(), 200, responseOfCreateEnvironment.Code)
			assert.Equal(suite.T(), saveEnvRequestDto.Environment, responseOfCreateEnvironment.Result.Environment)
		}

		appName := "app" + strings.ToLower(testUtils.GetRandomStringOfGivenLength(10))
		var teamId int
		if len(fetchAllTeamResponseDto.Result) > 0 {
			teamId = fetchAllTeamResponseDto.Result[0].Id
		} else {
			teamId = responseOfCreateProject.Result.Id
		}

		payLoadForCreateDevtronApp := CreatePayloadForDevtronApp(appName, teamId, 0, []PipelineConfigRouter.Labels{})
		byteValueOfStruct, _ := json.Marshal(payLoadForCreateDevtronApp)
		responseOfCreateDevtronApp := CreateDevtronApp(byteValueOfStruct, appName, teamId, 0, suite.authToken)
		assert.Equal(suite.T(), 200, responseOfCreateDevtronApp.Code)
		assert.Equal(suite.T(), appName, responseOfCreateDevtronApp.Result.AppName)
		assert.Equal(suite.T(), teamId, responseOfCreateDevtronApp.Result.TeamId)
		log.Println("=== App Name is :====>", responseOfCreateDevtronApp.Result.AppName)
		envConf := testUtils.ReadBaseEnvConfig()
		file := testUtils.ReadAnyJsonFile(envConf.ClassCredentialsFile)
		var configId int

		createAppApiResponse := responseOfCreateDevtronApp.Result
		workflowResponse := dtos.CreateWorkflowResponseDto{}
		if createAppApiResponsePtr != nil && workflowResponsePtr != nil {
			createAppApiResponse = *createAppApiResponsePtr
			workflowResponse = *workflowResponsePtr
		}
		log.Println("=== Here we are creating App Material ===")
		createAppMaterialRequestDto := PipelineConfigRouter.GetAppMaterialRequestDto(responseOfCreateDevtronApp.Result.Id, 1, false)
		appMaterialByteValue, _ := json.Marshal(createAppMaterialRequestDto)
		createAppMaterialResponse := PipelineConfigRouter.HitCreateAppMaterialApi(appMaterialByteValue, responseOfCreateDevtronApp.Result.Id, 1, false, suite.authToken)

		log.Println("=== Here we are saving docker build config ===")

		requestPayloadForSaveAppCiPipeline := PipelineConfigRouter.GetRequestPayloadForSaveAppCiPipeline(responseOfCreateDevtronApp.Result.Id, file.DockerRegistry, file.DockerUsername+"/test", file.DockerfilePath, file.DockerfileRepository, file.DockerfileRelativePath, createAppMaterialResponse.Result.Material[0].Id)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		PipelineConfigRouter.HitSaveAppCiPipeline(byteValueOfSaveAppCiPipeline, suite.authToken)

		log.Println("=== Here we are fetching latestChartReferenceId ===")
		time.Sleep(2 * time.Second)
		getChartReferenceResponse := PipelineConfigRouter.HitGetChartReferenceViaAppId(strconv.Itoa(responseOfCreateDevtronApp.Result.Id), suite.authToken)
		latestChartRef := getChartReferenceResponse.Result.LatestChartRef

		log.Println("=== Here we are fetching Template using getAppTemplateAPI ===")
		getTemplateResponse := PipelineConfigRouter.HitGetTemplateViaAppIdAndChartRefId(strconv.Itoa(responseOfCreateDevtronApp.Result.Id), strconv.Itoa(latestChartRef), suite.authToken)

		log.Println("=== Here we are fetching DefaultAppOverride from template response ===")
		defaultAppOverride := getTemplateResponse.Result.GlobalConfig.DefaultAppOverride

		log.Println("=== Here we are creating payload for SaveTemplate API ===")
		saveDeploymentTemplate := PipelineConfigRouter.GetRequestPayloadForSaveDeploymentTemplate(responseOfCreateDevtronApp.Result.Id, latestChartRef, defaultAppOverride)
		byteValueOfSaveDeploymentTemplate, _ := json.Marshal(saveDeploymentTemplate)
		jsonOfSaveDeploymentTemp := string(byteValueOfSaveDeploymentTemplate)
		jsonWithMicroserviceToleration, _ := sjson.Set(jsonOfSaveDeploymentTemp, "valuesOverride.tolerations.0", map[string]interface{}{"effect": "NoSchedule", "key": "microservice", "operator": "Equal", "value": "true"})
		finalJson, _ := sjson.Set(jsonWithMicroserviceToleration, "valuesOverride.tolerations.1", map[string]interface{}{"effect": "NoSchedule", "key": "kubernetes.azure.com/scalesetpriority", "operator": "Equal", "value": "spot"})
		updatedByteValueOfSaveDeploymentTemplate := []byte(finalJson)

		log.Println("=== Here we are hitting SaveTemplate API ===")
		PipelineConfigRouter.HitSaveDeploymentTemplateApi(updatedByteValueOfSaveDeploymentTemplate, suite.authToken)

		log.Println("=== Here we are saving Global Configmap ===")
		requestPayloadForConfigMap := HelperRouter.GetRequestPayloadForSecretOrConfig(0, "-config1", responseOfCreateDevtronApp.Result.Id, "environment", "kubernetes", false, false, false, false)
		byteValueOfSaverConfigMap, _ := json.Marshal(requestPayloadForConfigMap)
		globalConfigMap := HelperRouter.HitSaveGlobalConfigMap(byteValueOfSaverConfigMap, suite.authToken)
		configId = globalConfigMap.Result.Id

		log.Println("=== Here we are saving Global Secret ===")
		requestPayloadForSecret := HelperRouter.GetRequestPayloadForSecretOrConfig(configId, "-secret1", responseOfCreateDevtronApp.Result.Id, "environment", "kubernetes", false, false, true, false)
		byteValueOfSecret, _ := json.Marshal(requestPayloadForSecret)
		HelperRouter.HitSaveGlobalSecretApi(byteValueOfSecret, suite.authToken)

		workflowResponse = PipelineConfigRouter.HitCreateWorkflowApiWithFullPayload(responseOfCreateDevtronApp.Result.Id, suite.authToken, nil)

		log.Println("=== Here we are saving CD pipeline ===")
		payload := PipelineConfigRouter.GetRequestPayloadForSaveCdPipelineApi(responseOfCreateDevtronApp.Result.Id, workflowResponse.Result.AppWorkflowId, 1, workflowResponse.Result.CiPipelines[0].Id, workflowResponse.Result.CiPipelines[0].ParentCiPipeline, "AUTOMATIC", "", "", "AUTOMATIC")
		bytePayload, _ := json.Marshal(payload)
		savePipelineResponse := PipelineConfigRouter.HitSaveCdPipelineApi(bytePayload, suite.authToken)
		assert.Equal(suite.T(), 200, savePipelineResponse.Code)
		createAppApiResponsePtr = &createAppApiResponse
		workflowResponsePtr = &workflowResponse
		savePipelineResponsePtr = &savePipelineResponse
		log.Println("=== Here we are getting pipeline material ===")
		pipelineMaterial := PipelineConfigRouter.HitGetCiPipelineMaterial(workflowResponse.Result.CiPipelines[0].Id, suite.authToken)

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
		workflowStatus := PipelineConfigRouter.HitGetWorkflowStatus(responseOfCreateDevtronApp.Result.Id, suite.authToken)
		if workflowStatus.Result.CiWorkflowStatus[0].CiStatus == "Starting" {
			time.Sleep(5 * time.Second)
			workflowStatus = PipelineConfigRouter.HitGetWorkflowStatus(responseOfCreateDevtronApp.Result.Id, suite.authToken)
			assert.Equal(suite.T(), "Running", workflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		} else {
			assert.Equal(suite.T(), "Running", workflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		}
		log.Println("=== Here we are getting workflow and verifying the status after triggering via poll function ===")
		assert.True(suite.T(), PipelineConfigRouter.PollForGettingCdDeployStatusAfterTrigger(responseOfCreateDevtronApp.Result.Id, suite.authToken))
		updatedWorkflowStatus := PipelineConfigRouter.HitGetWorkflowStatus(responseOfCreateDevtronApp.Result.Id, suite.authToken)
		assert.Equal(suite.T(), "Succeeded", updatedWorkflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		assert.Equal(suite.T(), "Succeeded", updatedWorkflowStatus.Result.CdWorkflowStatus[0].DeployStatus)
		log.Println("=== Here we are Deleting the CD pipeline ===")
		deletePipelinePayload := PipelineConfigRouter.GetPayloadForDeleteCdPipeline(createAppApiResponsePtr.Id, savePipelineResponsePtr.Result.Pipelines[0].Id)
		deletePipelineByteCode, _ := json.Marshal(deletePipelinePayload)

		PipelineConfigRouter.HitForceDeleteCdPipelineApi(deletePipelineByteCode, suite.authToken)

		log.Println("=== Here we are Deleting the CI pipeline ===")
		PipelineConfigRouter.DeleteCiPipeline(createAppApiResponsePtr.Id, workflowResponsePtr.Result.CiPipelines[0].Id, suite.authToken)

		log.Println("=== Here we are Deleting CI Workflow ===")
		PipelineConfigRouter.HitDeleteWorkflowApi(createAppApiResponsePtr.Id, workflowResponsePtr.Result.AppWorkflowId, suite.authToken)
		log.Println("=== Here we are Deleting the app after all verifications ===")
		testUtils.DeleteApp(createAppApiResponsePtr.Id, createAppApiResponsePtr.AppName, createAppApiResponsePtr.TeamId, createAppApiResponsePtr.TemplateId, suite.authToken)

		if teamId == responseOfCreateProject.Result.Id {
			RbacFlows.DeleteProject(devtronDeletion.ProjectPayload, suite.authToken)
		}
		if len(allEnvironmentDetailsResponse.Result) == 0 {
			RbacFlows.DeleteEnv(devtronDeletion.EnvPayLoad, suite.authToken)
		}
	})

	suite.Run("A1=CreateDevtronAppWithGitRepoSubmodules", func() {
		var devtronDeletion RequestDTOs.RbacDevtronDeletion
		fetchAllTeamResponseDto := GetAllTeamsAutocomplete(suite.authToken)
		log.Println("Validating the response of FetchAllTeam API")
		var responseOfCreateProject TeamRouterResponseDTOs.SaveTeamResponseDTO
		if len(fetchAllTeamResponseDto.Result) == 0 {
			saveTeamRequestDto := TeamRouter.GetSaveTeamRequestDto()
			byteValueOfStruct, _ := json.Marshal(saveTeamRequestDto)
			responseOfCreateProject = RbacFlows.CreateProject(byteValueOfStruct, suite.authToken)
			devtronDeletion.ProjectPayload, _ = json.Marshal(responseOfCreateProject.Result)
			assert.Equal(suite.T(), 200, responseOfCreateProject.Code)
			assert.Equal(suite.T(), saveTeamRequestDto.Name, responseOfCreateProject.Result.Name)
		}
		queryParams := map[string]string{"auth": "true"}
		allEnvironmentDetailsResponse := GetAllEnvsAutocomplete(queryParams, suite.authToken)
		log.Println("Validating the response of GetAllEnvironmentDetails API")
		assert.NotNil(suite.T(), allEnvironmentDetailsResponse.Result)
		assert.Equal(suite.T(), 200, allEnvironmentDetailsResponse.Code)
		if len(allEnvironmentDetailsResponse.Result) == 0 {
			saveEnvRequestDto := RbacFlows.GetSaveEnvRequestDto()
			saveEnvRequestDto.EnvironmentIdentifier = "default_cluster__" + saveEnvRequestDto.Namespace
			byteValueOfStruct, _ := json.Marshal(saveEnvRequestDto)
			responseOfCreateEnvironment := RbacFlows.CreateEnv(byteValueOfStruct, suite.authToken)
			devtronDeletion.EnvPayLoad, _ = json.Marshal(responseOfCreateEnvironment.Result)
			assert.Equal(suite.T(), 200, responseOfCreateEnvironment.Code)
			assert.Equal(suite.T(), saveEnvRequestDto.Environment, responseOfCreateEnvironment.Result.Environment)
		}

		appName := "app" + strings.ToLower(testUtils.GetRandomStringOfGivenLength(10))
		var teamId int
		if len(fetchAllTeamResponseDto.Result) > 0 {
			teamId = fetchAllTeamResponseDto.Result[0].Id
		} else {
			teamId = responseOfCreateProject.Result.Id
		}

		payLoadForCreateDevtronApp := CreatePayloadForDevtronApp(appName, teamId, 0, []PipelineConfigRouter.Labels{})
		byteValueOfStruct, _ := json.Marshal(payLoadForCreateDevtronApp)
		responseOfCreateDevtronApp := CreateDevtronApp(byteValueOfStruct, appName, teamId, 0, suite.authToken)
		assert.Equal(suite.T(), 200, responseOfCreateDevtronApp.Code)
		assert.Equal(suite.T(), appName, responseOfCreateDevtronApp.Result.AppName)
		assert.Equal(suite.T(), teamId, responseOfCreateDevtronApp.Result.TeamId)
		log.Println("=== App Name is :====>", responseOfCreateDevtronApp.Result.AppName)
		envConf := testUtils.ReadBaseEnvConfig()
		file := testUtils.ReadAnyJsonFile(envConf.ClassCredentialsFile)
		var configId int

		createAppApiResponse := responseOfCreateDevtronApp.Result
		workflowResponse := dtos.CreateWorkflowResponseDto{}
		if createAppApiResponsePtr != nil && workflowResponsePtr != nil {
			createAppApiResponse = *createAppApiResponsePtr
			workflowResponse = *workflowResponsePtr
		}
		log.Println("=== Here we are creating App Material ===")
		createAppMaterialRequestDto := PipelineConfigRouter.GetAppMaterialRequestDto(responseOfCreateDevtronApp.Result.Id, 1, true)
		createAppMaterialRequestDto.Materials[0].Url = SUBMODULES_GIT_REPO_URL
		appMaterialByteValue, _ := json.Marshal(createAppMaterialRequestDto)
		createAppMaterialResponse := PipelineConfigRouter.HitCreateAppMaterialApi(appMaterialByteValue, responseOfCreateDevtronApp.Result.Id, 1, true, suite.authToken)

		log.Println("=== Here we are saving docker build config ===")

		requestPayloadForSaveAppCiPipeline := PipelineConfigRouter.GetRequestPayloadForSaveAppCiPipeline(responseOfCreateDevtronApp.Result.Id, file.DockerRegistry, file.DockerUsername+"/test", file.DockerfilePath, file.DockerfileRepository, file.DockerfileRelativePath, createAppMaterialResponse.Result.Material[0].Id)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		PipelineConfigRouter.HitSaveAppCiPipeline(byteValueOfSaveAppCiPipeline, suite.authToken)

		log.Println("=== Here we are fetching latestChartReferenceId ===")
		time.Sleep(2 * time.Second)
		getChartReferenceResponse := PipelineConfigRouter.HitGetChartReferenceViaAppId(strconv.Itoa(responseOfCreateDevtronApp.Result.Id), suite.authToken)
		latestChartRef := getChartReferenceResponse.Result.LatestChartRef

		log.Println("=== Here we are fetching Template using getAppTemplateAPI ===")
		getTemplateResponse := PipelineConfigRouter.HitGetTemplateViaAppIdAndChartRefId(strconv.Itoa(responseOfCreateDevtronApp.Result.Id), strconv.Itoa(latestChartRef), suite.authToken)

		log.Println("=== Here we are fetching DefaultAppOverride from template response ===")
		defaultAppOverride := getTemplateResponse.Result.GlobalConfig.DefaultAppOverride

		log.Println("=== Here we are creating payload for SaveTemplate API ===")
		saveDeploymentTemplate := PipelineConfigRouter.GetRequestPayloadForSaveDeploymentTemplate(responseOfCreateDevtronApp.Result.Id, latestChartRef, defaultAppOverride)
		byteValueOfSaveDeploymentTemplate, _ := json.Marshal(saveDeploymentTemplate)
		jsonOfSaveDeploymentTemp := string(byteValueOfSaveDeploymentTemplate)
		jsonWithMicroserviceToleration, _ := sjson.Set(jsonOfSaveDeploymentTemp, "valuesOverride.tolerations.0", map[string]interface{}{"effect": "NoSchedule", "key": "microservice", "operator": "Equal", "value": "true"})
		finalJson, _ := sjson.Set(jsonWithMicroserviceToleration, "valuesOverride.tolerations.1", map[string]interface{}{"effect": "NoSchedule", "key": "kubernetes.azure.com/scalesetpriority", "operator": "Equal", "value": "spot"})
		updatedByteValueOfSaveDeploymentTemplate := []byte(finalJson)

		log.Println("=== Here we are hitting SaveTemplate API ===")
		PipelineConfigRouter.HitSaveDeploymentTemplateApi(updatedByteValueOfSaveDeploymentTemplate, suite.authToken)

		log.Println("=== Here we are saving Global Configmap ===")
		requestPayloadForConfigMap := HelperRouter.GetRequestPayloadForSecretOrConfig(0, "-config1", responseOfCreateDevtronApp.Result.Id, "environment", "kubernetes", false, false, false, false)
		byteValueOfSaverConfigMap, _ := json.Marshal(requestPayloadForConfigMap)
		globalConfigMap := HelperRouter.HitSaveGlobalConfigMap(byteValueOfSaverConfigMap, suite.authToken)
		configId = globalConfigMap.Result.Id

		log.Println("=== Here we are saving Global Secret ===")
		requestPayloadForSecret := HelperRouter.GetRequestPayloadForSecretOrConfig(configId, "-secret1", responseOfCreateDevtronApp.Result.Id, "environment", "kubernetes", false, false, true, false)
		byteValueOfSecret, _ := json.Marshal(requestPayloadForSecret)
		HelperRouter.HitSaveGlobalSecretApi(byteValueOfSecret, suite.authToken)

		workflowResponse = PipelineConfigRouter.HitCreateWorkflowApiWithFullPayload(responseOfCreateDevtronApp.Result.Id, suite.authToken, []string{"test"})

		log.Println("=== Here we are saving CD pipeline ===")
		payload := PipelineConfigRouter.GetRequestPayloadForSaveCdPipelineApi(responseOfCreateDevtronApp.Result.Id, workflowResponse.Result.AppWorkflowId, 1, workflowResponse.Result.CiPipelines[0].Id, workflowResponse.Result.CiPipelines[0].ParentCiPipeline, "AUTOMATIC", "", "", "AUTOMATIC")
		bytePayload, _ := json.Marshal(payload)
		savePipelineResponse := PipelineConfigRouter.HitSaveCdPipelineApi(bytePayload, suite.authToken)
		assert.Equal(suite.T(), 200, savePipelineResponse.Code)
		createAppApiResponsePtr = &createAppApiResponse
		workflowResponsePtr = &workflowResponse
		savePipelineResponsePtr = &savePipelineResponse
		log.Println("=== Here we are getting pipeline material ===")
		pipelineMaterial := PipelineConfigRouter.HitGetCiPipelineMaterial(workflowResponse.Result.CiPipelines[0].Id, suite.authToken)

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
		workflowStatus := PipelineConfigRouter.HitGetWorkflowStatus(responseOfCreateDevtronApp.Result.Id, suite.authToken)
		if workflowStatus.Result.CiWorkflowStatus[0].CiStatus == "Starting" {
			time.Sleep(5 * time.Second)
			workflowStatus = PipelineConfigRouter.HitGetWorkflowStatus(responseOfCreateDevtronApp.Result.Id, suite.authToken)
			assert.Equal(suite.T(), "Running", workflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		} else {
			assert.Equal(suite.T(), "Running", workflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		}
		log.Println("=== Here we are getting workflow and verifying the status after triggering via poll function ===")
		assert.True(suite.T(), PipelineConfigRouter.PollForGettingCdDeployStatusAfterTrigger(responseOfCreateDevtronApp.Result.Id, suite.authToken))
		updatedWorkflowStatus := PipelineConfigRouter.HitGetWorkflowStatus(responseOfCreateDevtronApp.Result.Id, suite.authToken)
		assert.Equal(suite.T(), "Succeeded", updatedWorkflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		assert.Equal(suite.T(), "Succeeded", updatedWorkflowStatus.Result.CdWorkflowStatus[0].DeployStatus)
		log.Println("=== Here we are Deleting the CD pipeline ===")
		deletePipelinePayload := PipelineConfigRouter.GetPayloadForDeleteCdPipeline(createAppApiResponsePtr.Id, savePipelineResponsePtr.Result.Pipelines[0].Id)
		deletePipelineByteCode, _ := json.Marshal(deletePipelinePayload)

		PipelineConfigRouter.HitForceDeleteCdPipelineApi(deletePipelineByteCode, suite.authToken)

		log.Println("=== Here we are Deleting the CI pipeline ===")
		PipelineConfigRouter.DeleteCiPipeline(createAppApiResponsePtr.Id, workflowResponsePtr.Result.CiPipelines[0].Id, suite.authToken)

		log.Println("=== Here we are Deleting CI Workflow ===")
		PipelineConfigRouter.HitDeleteWorkflowApi(createAppApiResponsePtr.Id, workflowResponsePtr.Result.AppWorkflowId, suite.authToken)
		log.Println("=== Here we are Deleting the app after all verifications ===")
		testUtils.DeleteApp(createAppApiResponsePtr.Id, createAppApiResponsePtr.AppName, createAppApiResponsePtr.TeamId, createAppApiResponsePtr.TemplateId, suite.authToken)

		if teamId == responseOfCreateProject.Result.Id {
			RbacFlows.DeleteProject(devtronDeletion.ProjectPayload, suite.authToken)
		}
		if len(allEnvironmentDetailsResponse.Result) == 0 {
			RbacFlows.DeleteEnv(devtronDeletion.EnvPayLoad, suite.authToken)
		}
	})
}
