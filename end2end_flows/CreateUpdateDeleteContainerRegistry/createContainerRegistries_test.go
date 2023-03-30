package CreateUpdateDeleteContainerRegistry

import (
	"automation-suite/PipelineConfigRouter"
	"automation-suite/dockerRegRouter"
	"automation-suite/dockerRegRouter/ResponseDTOs"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"time"
)

func (suite *CreateContainerRegistryFlowsTestSuite) Test() {
	var byteValueOfSaveDockerRegistry []byte
	var saveDockerRegistryResponseDto ResponseDTOs.SaveDockerRegistryResponseDto
	var getDockerRegistryResponseDto ResponseDTOs.SaveDockerRegistryResponseDto
	var getAllDockerRegistryResponseDto ResponseDTOs.GetAllDockerRegistryResponseDto
	suite.Run("HitApiWithValidCredentials", func() {
		getAllDockerRegistryResponseDto = dockerRegRouter.HitGetAllDockerRegistryApi(suite.authToken)
		saveDockerRegistryRequestDto := dockerRegRouter.GetDockerRegistryRequestDto(false)
		byteValueOfSaveDockerRegistry, _ = json.Marshal(saveDockerRegistryRequestDto)

		log.Println("=== Hitting The post Docker registry API ===")
		saveDockerRegistryResponseDto = dockerRegRouter.HitSaveContainerRegistryApi(byteValueOfSaveDockerRegistry, suite.authToken)

		newGetAllDockerRegistryResponseDto := dockerRegRouter.HitGetAllDockerRegistryApi(suite.authToken)
		assert.Equal(suite.T(), len(newGetAllDockerRegistryResponseDto.Result), len(getAllDockerRegistryResponseDto.Result)+1)
		log.Println("=== Validating the Response of the save docker registry API... === ")
		assert.Equal(suite.T(), saveDockerRegistryRequestDto.Id, saveDockerRegistryResponseDto.Result.Id)
		assert.Equal(suite.T(), saveDockerRegistryRequestDto.IsDefault, saveDockerRegistryResponseDto.Result.IsDefault)

		getDockerRegistryResponseDto = dockerRegRouter.HitGetDockerRegistryApi(suite.authToken, saveDockerRegistryResponseDto.Result.Id)
		assert.Equal(suite.T(), getDockerRegistryResponseDto.Result.Id, saveDockerRegistryResponseDto.Result.Id)

		envConf := Base.ReadBaseEnvConfig()
		config := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)

		log.Println("=== Creating a App ===")
		createAppApiResponse := Base.CreateApp(suite.authToken).Result

		log.Println("=== Creating App Material ===")
		createAppMaterialRequestDto := PipelineConfigRouter.GetAppMaterialRequestDto(createAppApiResponse.Id, 1, false)
		appMaterialByteValue, _ := json.Marshal(createAppMaterialRequestDto)
		createAppMaterialResponse := PipelineConfigRouter.HitCreateAppMaterialApi(appMaterialByteValue, createAppApiResponse.Id, 1, false, suite.authToken)

		log.Println("=== Saving docker build config ===")
		requestPayloadForSaveAppCiPipeline := PipelineConfigRouter.GetRequestPayloadForSaveAppCiPipeline(createAppApiResponse.Id, saveDockerRegistryResponseDto.Result.Id, saveDockerRegistryResponseDto.Result.Id+"/test", config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, createAppMaterialResponse.Result.Material[0].Id)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		PipelineConfigRouter.HitSaveAppCiPipeline(byteValueOfSaveAppCiPipeline, suite.authToken)

		log.Println("=== Fetching latestChartReferenceId ===")
		getChartReferenceResponse := PipelineConfigRouter.HitGetChartReferenceViaAppId(strconv.Itoa(createAppApiResponse.Id), suite.authToken)
		latestChartRef := getChartReferenceResponse.Result.LatestChartRef

		log.Println("=== Fetching DefaultAppOverride using getAppTemplateAPI ===")
		getTemplateResponse := PipelineConfigRouter.HitGetTemplateViaAppIdAndChartRefId(strconv.Itoa(createAppApiResponse.Id), strconv.Itoa(latestChartRef), suite.authToken)

		log.Println("=== Fetching DefaultAppOverride using getAppTemplateAPI ===")
		defaultAppOverride := getTemplateResponse.Result.GlobalConfig.DefaultAppOverride

		log.Println("=== Creating payload for SaveTemplate API ===")
		saveDeploymentTemplate := PipelineConfigRouter.GetRequestPayloadForSaveDeploymentTemplate(createAppApiResponse.Id, latestChartRef, defaultAppOverride)
		byteValueOfSaveDeploymentTemplate, _ := json.Marshal(saveDeploymentTemplate)

		log.Println("=== Hitting SaveTemplate API ===")
		PipelineConfigRouter.HitSaveDeploymentTemplateApi(byteValueOfSaveDeploymentTemplate, suite.authToken)

		log.Println("=== Hitting CreateWorkFlow API ===")
		workflowResponse := PipelineConfigRouter.HitCreateWorkflowApiWithFullPayload(createAppApiResponse.Id, suite.authToken).Result

		log.Println("=== Hitting GetCiPipelineMaterial API ===")
		pipelineMaterial := PipelineConfigRouter.HitGetCiPipelineMaterial(workflowResponse.CiPipelines[0].Id, suite.authToken)
		payloadForTriggerCiPipeline := PipelineConfigRouter.CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.CiPipelines[0].Id, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)

		log.Println("=== Hitting TriggerCiPipeline API ===")
		triggerCiPipelineResponse := PipelineConfigRouter.HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		time.After(60 * time.Second)

		log.Println("=== Reading Logs ===")
		suite.checkForCiLogs(strconv.Itoa(workflowResponse.CiPipelines[0].Id), triggerCiPipelineResponse.Result.ApiResponse, 151, "Login Succeeded")

		log.Println("=== Here we are Deleting the CI pipeline ===")
		PipelineConfigRouter.DeleteCiPipeline(createAppApiResponse.Id, workflowResponse.CiPipelines[0].Id, suite.authToken)

		log.Println("=== Here we are Deleting CI Workflow ===")
		PipelineConfigRouter.HitDeleteWorkflowApi(createAppApiResponse.Id, workflowResponse.AppWorkflowId, suite.authToken)

		log.Println("=== Here we are Deleting the app after all verifications ===")
		Base.DeleteApp(createAppApiResponse.Id, createAppApiResponse.AppName, createAppApiResponse.TeamId, createAppApiResponse.TemplateId, suite.authToken)

		log.Println("=== getting payload for Delete registry API ===")
		byteValueOfDeleteDockerRegistry := dockerRegRouter.GetPayLoadForDeleteDockerRegistryAPI(saveDockerRegistryResponseDto.Result.Id, saveDockerRegistryResponseDto.Result.IpsConfig.Id, saveDockerRegistryResponseDto.Result.PluginId, saveDockerRegistryResponseDto.Result.RegistryUrl, saveDockerRegistryResponseDto.Result.RegistryType, saveDockerRegistryResponseDto.Result.Username, saveDockerRegistryResponseDto.Result.Password, saveDockerRegistryResponseDto.Result.IsDefault)

		log.Println("=== Hitting the Delete team API for Removing the data created via automation ===")
		dockerRegRouter.HitDeleteDockerRegistryApi(byteValueOfDeleteDockerRegistry, suite.authToken)
	})

	suite.Run("HitApiWithInvalidCredentials", func() {
		getAllDockerRegistryResponseDto = dockerRegRouter.HitGetAllDockerRegistryApi(suite.authToken)
		saveDockerRegistryRequestDto := dockerRegRouter.GetDockerRegistryRequestDto(false)
		saveDockerRegistryRequestDto.Password = "123"
		byteValueOfSaveDockerRegistry, _ = json.Marshal(saveDockerRegistryRequestDto)

		log.Println("Hitting The post Docker registry API")
		saveDockerRegistryResponseDto = dockerRegRouter.HitSaveContainerRegistryApi(byteValueOfSaveDockerRegistry, suite.authToken)

		newGetAllDockerRegistryResponseDto := dockerRegRouter.HitGetAllDockerRegistryApi(suite.authToken)
		assert.Equal(suite.T(), len(newGetAllDockerRegistryResponseDto.Result), len(getAllDockerRegistryResponseDto.Result)+1)
		log.Println("Validating the Response of the save docker registry API...")
		assert.Equal(suite.T(), saveDockerRegistryRequestDto.Id, saveDockerRegistryResponseDto.Result.Id)
		assert.Equal(suite.T(), saveDockerRegistryRequestDto.IsDefault, saveDockerRegistryResponseDto.Result.IsDefault)

		getDockerRegistryResponseDto = dockerRegRouter.HitGetDockerRegistryApi(suite.authToken, saveDockerRegistryResponseDto.Result.Id)
		assert.Equal(suite.T(), getDockerRegistryResponseDto.Result.Id, saveDockerRegistryResponseDto.Result.Id)

		envConf := Base.ReadBaseEnvConfig()
		config := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)

		log.Println("=== Creating a App ===")
		createAppApiResponse := Base.CreateApp(suite.authToken).Result

		log.Println("=== Creating App Material ===")
		createAppMaterialRequestDto := PipelineConfigRouter.GetAppMaterialRequestDto(createAppApiResponse.Id, 1, false)
		appMaterialByteValue, _ := json.Marshal(createAppMaterialRequestDto)
		createAppMaterialResponse := PipelineConfigRouter.HitCreateAppMaterialApi(appMaterialByteValue, createAppApiResponse.Id, 1, false, suite.authToken)

		log.Println("=== Saving docker build config ===")
		requestPayloadForSaveAppCiPipeline := PipelineConfigRouter.GetRequestPayloadForSaveAppCiPipeline(createAppApiResponse.Id, saveDockerRegistryResponseDto.Result.Id, saveDockerRegistryResponseDto.Result.Id+"/test", config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, createAppMaterialResponse.Result.Material[0].Id)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		PipelineConfigRouter.HitSaveAppCiPipeline(byteValueOfSaveAppCiPipeline, suite.authToken)

		log.Println("=== Fetching latestChartReferenceId ===")
		getChartReferenceResponse := PipelineConfigRouter.HitGetChartReferenceViaAppId(strconv.Itoa(createAppApiResponse.Id), suite.authToken)
		latestChartRef := getChartReferenceResponse.Result.LatestChartRef

		log.Println("=== Fetching DefaultAppOverride using getAppTemplateAPI ===")
		getTemplateResponse := PipelineConfigRouter.HitGetTemplateViaAppIdAndChartRefId(strconv.Itoa(createAppApiResponse.Id), strconv.Itoa(latestChartRef), suite.authToken)

		log.Println("=== Fetching DefaultAppOverride using getAppTemplateAPI ===")
		defaultAppOverride := getTemplateResponse.Result.GlobalConfig.DefaultAppOverride

		log.Println("=== Creating payload for SaveTemplate API ===")
		saveDeploymentTemplate := PipelineConfigRouter.GetRequestPayloadForSaveDeploymentTemplate(createAppApiResponse.Id, latestChartRef, defaultAppOverride)
		byteValueOfSaveDeploymentTemplate, _ := json.Marshal(saveDeploymentTemplate)

		log.Println("=== Hitting SaveTemplate API ===")
		PipelineConfigRouter.HitSaveDeploymentTemplateApi(byteValueOfSaveDeploymentTemplate, suite.authToken)
		log.Println("=== Hitting CreateWorkFlow API ===")
		workflowResponse := PipelineConfigRouter.HitCreateWorkflowApiWithFullPayload(createAppApiResponse.Id, suite.authToken).Result
		log.Println("=== Hitting GetCiPipelineMaterial API ===")
		pipelineMaterial := PipelineConfigRouter.HitGetCiPipelineMaterial(workflowResponse.CiPipelines[0].Id, suite.authToken)
		payloadForTriggerCiPipeline := PipelineConfigRouter.CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.CiPipelines[0].Id, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		log.Println("=== Hitting TriggerCiPipeline API ===")
		triggerCiPipelineResponse := PipelineConfigRouter.HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		log.Println("=== Reading Logs ===")
		suite.checkForCiLogs(strconv.Itoa(workflowResponse.CiPipelines[0].Id), triggerCiPipelineResponse.Result.ApiResponse, 151, "Login Succeeded")
		log.Println("=== Here we are Deleting the CI pipeline ===")
		PipelineConfigRouter.DeleteCiPipeline(createAppApiResponse.Id, workflowResponse.CiPipelines[0].Id, suite.authToken)

		log.Println("=== Here we are Deleting CI Workflow ===")
		PipelineConfigRouter.HitDeleteWorkflowApi(createAppApiResponse.Id, workflowResponse.AppWorkflowId, suite.authToken)

		log.Println("=== Here we are Deleting the app after all verifications ===")
		Base.DeleteApp(createAppApiResponse.Id, createAppApiResponse.AppName, createAppApiResponse.TeamId, createAppApiResponse.TemplateId, suite.authToken)
		log.Println("getting payload for Delete registry API")
		byteValueOfDeleteDockerRegistry := dockerRegRouter.GetPayLoadForDeleteDockerRegistryAPI(saveDockerRegistryResponseDto.Result.Id, saveDockerRegistryResponseDto.Result.IpsConfig.Id, saveDockerRegistryResponseDto.Result.PluginId, saveDockerRegistryResponseDto.Result.RegistryUrl, saveDockerRegistryResponseDto.Result.RegistryType, saveDockerRegistryResponseDto.Result.Username, saveDockerRegistryResponseDto.Result.Password, saveDockerRegistryResponseDto.Result.IsDefault)
		log.Println("Hitting the Delete team API for Removing the data created via automation")
		dockerRegRouter.HitDeleteDockerRegistryApi(byteValueOfDeleteDockerRegistry, suite.authToken)
	})
}
