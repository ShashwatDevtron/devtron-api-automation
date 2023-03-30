package GitProviderFullFlow

import (
	"automation-suite/GitHostRouter"
	"automation-suite/GitProviderRouter/ResponseDTOs"
	"automation-suite/PipelineConfigRouter"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
)

func (suite *GitProRouterTestSuite) TestClassGetGitProivder() {

	//var byteValueOfSaveGitProvider []byte
	var saveGitProviderResponseDto ResponseDTOs.SaveGitProviderResponseDto
	var lengthOfGitProviders int

	suite.Run("Checking Full Flow of Git Provider with valid payload", func() {

		log.Println("Hitting The get git provider API")
		getGitProviderResponse := HitGetListOfGitProviders(suite.authToken)

		lengthOfGitProviders = len(getGitProviderResponse.Result)

		saveGitHostListResponse := GitHostRouter.HitGetGitHostApi(suite.authToken)

		log.Println("SaveGitProviderWithValidPayload")
		saveGitProviderRequestDto := GetGitProviderRequestObjectDto(saveGitHostListResponse.Result[0].Id)
		byteValueOfSaveGitProvider, _ := json.Marshal(saveGitProviderRequestDto)
		log.Println("Hitting The post git provider API")
		saveGitProviderResponseDto = HitSaveOneGitProviderApi(byteValueOfSaveGitProvider, suite.authToken)
		log.Println("Validating the Response of the save git provider API...")
		assert.Equal(suite.T(), saveGitProviderRequestDto.Name, saveGitProviderResponseDto.Result.Name)

		log.Println("Hitting The get git provider API")
		getGitProviderResponses := HitGetListOfGitProviders(suite.authToken)
		assert.Equal(suite.T(), len(getGitProviderResponses.Result), lengthOfGitProviders+1)

		envConf := Base.ReadBaseEnvConfig()
		config := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)

		log.Println("=== Here we are creating a App ===")
		createAppApiResponse := CreateAppWithRandomData(suite.authToken).Result

		log.Println("=== Here we are creating App Material ===")
		createAppMaterialResponse := HitCreateAppMaterialWithGitProviderIdApi(nil, createAppApiResponse.Id, saveGitProviderResponseDto.Result.Id, false, suite.authToken)

		log.Println("=== Saving docker build config ===")
		requestPayloadForSaveAppCiPipeline := GetRequestPayloadForSavingAppCiPipeline(createAppApiResponse.Id, config.DockerRegistry, config.DockerRegistry+"/test", config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, createAppMaterialResponse.Result.Material[0].Id)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		HitSaveCiPipelineForApp(byteValueOfSaveAppCiPipeline, suite.authToken)

		log.Println("=== Fetching latestChartReferenceId ===")
		getChartReferenceResponse := HitGetChartRefViaAppId(strconv.Itoa(createAppApiResponse.Id), suite.authToken)
		latestChartRef := getChartReferenceResponse.Result.LatestChartRef

		log.Println("=== Fetching DefaultAppOverride using getAppTemplateAPI ===")
		getTemplateResponse := HitGetTemplateViaAppIdAndChartRef(strconv.Itoa(createAppApiResponse.Id), strconv.Itoa(latestChartRef), suite.authToken)

		log.Println("=== Fetching DefaultAppOverride using getAppTemplateAPI ===")
		defaultAppOverride := getTemplateResponse.Result.GlobalConfig.DefaultAppOverride

		log.Println("=== Creating payload for SaveTemplate API ===")
		saveDeploymentTemplate := GetRequestPayloadForSavingDeploymentTemplate(createAppApiResponse.Id, latestChartRef, defaultAppOverride)
		byteValueOfSaveDeploymentTemplate, _ := json.Marshal(saveDeploymentTemplate)

		log.Println("=== Hitting SaveTemplate API ===")
		HitSavingDeploymentTemplateApi(byteValueOfSaveDeploymentTemplate, suite.authToken)

		workflowResponse := HitCreatingWorkflowApiWithFullPayload(createAppApiResponse.Id, suite.authToken).Result
		log.Println("=== Hitting GetCiPipelineMaterial API ===")
		pipelineMaterial := HitGettingCiPipelineMaterial(workflowResponse.CiPipelines[0].Id, suite.authToken)
		payloadForTriggerCiPipeline := CreatePayloadForTriggeringCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.CiPipelines[0].Id, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		log.Println("=== Hitting TriggerCiPipeline API ===")
		triggerCiPipelineResponse := HittingTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)

		suite.checkForCiLogs(strconv.Itoa(workflowResponse.CiPipelines[0].Id), triggerCiPipelineResponse.Result.ApiResponse, 2, 3)
		log.Println("=== Here we are Deleting the CI pipeline ===")
		PipelineConfigRouter.DeleteCiPipeline(createAppApiResponse.Id, workflowResponse.CiPipelines[0].Id, suite.authToken)

		log.Println("=== Here we are Deleting CI Workflow ===")
		PipelineConfigRouter.HitDeleteWorkflowApi(createAppApiResponse.Id, workflowResponse.AppWorkflowId, suite.authToken)

		log.Println("=== Here we are Deleting the app after all verifications ===")
		Base.DeleteApp(createAppApiResponse.Id, createAppApiResponse.AppName, createAppApiResponse.TeamId, createAppApiResponse.TemplateId, suite.authToken)

		log.Println("getting payload for Update git provider API")
		byteValueOfUpdateGitProvider := GetPayLoadForUpdateOneGitProviderAPI(saveGitProviderResponseDto.Result.Id, saveGitProviderResponseDto.Result.GitHostId)
		log.Println("Hitting the Update Git API for Updating the data created via automation")
		HitUpdateOneGitProviderApiResponse := HitUpdateOneGitProviderApi(byteValueOfUpdateGitProvider, suite.authToken)
		assert.NotEqual(suite.T(), saveGitProviderResponseDto.Result.Name, HitUpdateOneGitProviderApiResponse.Result.Name)
		assert.Equal(suite.T(), saveGitProviderResponseDto.Result.Id, HitUpdateOneGitProviderApiResponse.Result.Id)

		log.Println("getting payload for Delete git provider API")
		byteValueOfDeleteDockerRegistry := GetPayLoadForDeleteOneGitProviderAPI(saveGitProviderResponseDto.Result.Id, saveGitProviderResponseDto.Result.GitHostId, saveGitProviderResponseDto.Result.Url, saveGitProviderResponseDto.Result.AuthMode, saveGitProviderResponseDto.Result.Name)
		log.Println("Hitting the Delete Git API for Removing the data created via automation")
		HitDeleteOneGitProviderApi(byteValueOfDeleteDockerRegistry, suite.authToken)

		log.Println("Hitting The get git provider API")
		getGitProviderResponseses := HitGetListOfGitProviders(suite.authToken)
		assert.Equal(suite.T(), len(getGitProviderResponseses.Result), lengthOfGitProviders)
	})

	suite.Run("Checking Full Flow of Git Provider with Invalid GitUrl", func() {

		log.Println("Hitting The get git provider API")
		getGitProviderResponse := HitGetListOfGitProviders(suite.authToken)

		lengthOfGitProviders = len(getGitProviderResponse.Result)

		saveGitHostListResponse := GitHostRouter.HitGetGitHostApi(suite.authToken)

		log.Println("SaveGitProviderWithValidPayload")
		saveGitProviderRequestDto := GetGitProviderRequestObjectDto(saveGitHostListResponse.Result[0].Id)
		byteValueOfSaveGitProvider, _ := json.Marshal(saveGitProviderRequestDto)
		log.Println("Hitting The post git provider API")
		saveGitProviderResponseDto = HitSaveOneGitProviderApi(byteValueOfSaveGitProvider, suite.authToken)
		log.Println("Validating the Response of the save git provider API...")
		assert.Equal(suite.T(), saveGitProviderRequestDto.Name, saveGitProviderResponseDto.Result.Name)

		log.Println("Hitting The get git provider API")
		getGitProviderResponses := HitGetListOfGitProviders(suite.authToken)
		assert.Equal(suite.T(), len(getGitProviderResponses.Result), lengthOfGitProviders+1)

		envConf := Base.ReadBaseEnvConfig()
		config := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)

		log.Println("=== Here we are creating a App ===")
		createAppApiResponse := CreateAppWithRandomData(suite.authToken).Result

		log.Println("=== Here we are creating App Material ===")
		createAppMaterialResponse := HitCreateAppMaterialWithGitProviderIdApi(nil, createAppApiResponse.Id, saveGitProviderResponseDto.Result.Id, false, suite.authToken)

		log.Println("=== Saving docker build config ===")
		requestPayloadForSaveAppCiPipeline := GetRequestPayloadForSavingAppCiPipeline(createAppApiResponse.Id, config.DockerRegistry, config.DockerRegistry+"/test", config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, createAppMaterialResponse.Result.Material[0].Id)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		HitSaveCiPipelineForApp(byteValueOfSaveAppCiPipeline, suite.authToken)

		log.Println("=== Fetching latestChartReferenceId ===")
		getChartReferenceResponse := HitGetChartRefViaAppId(strconv.Itoa(createAppApiResponse.Id), suite.authToken)
		latestChartRef := getChartReferenceResponse.Result.LatestChartRef

		log.Println("=== Fetching DefaultAppOverride using getAppTemplateAPI ===")
		getTemplateResponse := HitGetTemplateViaAppIdAndChartRef(strconv.Itoa(createAppApiResponse.Id), strconv.Itoa(latestChartRef), suite.authToken)

		log.Println("=== Fetching DefaultAppOverride using getAppTemplateAPI ===")
		defaultAppOverride := getTemplateResponse.Result.GlobalConfig.DefaultAppOverride

		log.Println("=== Creating payload for SaveTemplate API ===")
		saveDeploymentTemplate := GetRequestPayloadForSavingDeploymentTemplate(createAppApiResponse.Id, latestChartRef, defaultAppOverride)
		byteValueOfSaveDeploymentTemplate, _ := json.Marshal(saveDeploymentTemplate)

		log.Println("=== Hitting SaveTemplate API ===")
		HitSavingDeploymentTemplateApi(byteValueOfSaveDeploymentTemplate, suite.authToken)

		workflowResponse := HitCreatingWorkflowApiWithFullPayload(createAppApiResponse.Id, suite.authToken).Result
		log.Println("=== Hitting GetCiPipelineMaterial API ===")
		pipelineMaterial := HitGettingCiPipelineMaterial(workflowResponse.CiPipelines[0].Id, suite.authToken)
		//payloadForTriggerCiPipeline := CreatePayloadForTriggeringCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.CiPipelines[0].Id, pipelineMaterial.Result[0].Id, true)
		assert.Equal(suite.T(), 0, len(pipelineMaterial.Result[0].History))

		log.Println("=== Here we are Deleting the CI pipeline ===")
		PipelineConfigRouter.DeleteCiPipeline(createAppApiResponse.Id, workflowResponse.CiPipelines[0].Id, suite.authToken)

		log.Println("=== Here we are Deleting CI Workflow ===")
		PipelineConfigRouter.HitDeleteWorkflowApi(createAppApiResponse.Id, workflowResponse.AppWorkflowId, suite.authToken)

		log.Println("=== Here we are Deleting the app after all verifications ===")
		Base.DeleteApp(createAppApiResponse.Id, createAppApiResponse.AppName, createAppApiResponse.TeamId, createAppApiResponse.TemplateId, suite.authToken)

		log.Println("getting payload for Update git provider API")
		byteValueOfUpdateGitProvider := GetPayLoadForUpdateOneGitProviderAPI(saveGitProviderResponseDto.Result.Id, saveGitProviderResponseDto.Result.GitHostId)
		log.Println("Hitting the Update Git API for Updating the data created via automation")
		HitUpdateOneGitProviderApiResponse := HitUpdateOneGitProviderApi(byteValueOfUpdateGitProvider, suite.authToken)
		assert.NotEqual(suite.T(), saveGitProviderResponseDto.Result.Name, HitUpdateOneGitProviderApiResponse.Result.Name)
		assert.Equal(suite.T(), saveGitProviderResponseDto.Result.Id, HitUpdateOneGitProviderApiResponse.Result.Id)

		log.Println("getting payload for Delete git provider API")
		byteValueOfDeleteDockerRegistry := GetPayLoadForDeleteOneGitProviderAPI(saveGitProviderResponseDto.Result.Id, saveGitProviderResponseDto.Result.GitHostId, saveGitProviderResponseDto.Result.Url, saveGitProviderResponseDto.Result.AuthMode, saveGitProviderResponseDto.Result.Name)
		log.Println("Hitting the Delete Git API for Removing the data created via automation")
		HitDeleteOneGitProviderApi(byteValueOfDeleteDockerRegistry, suite.authToken)

		log.Println("Hitting The get git provider API")
		getGitProviderResponseses := HitGetListOfGitProviders(suite.authToken)
		assert.Equal(suite.T(), len(getGitProviderResponseses.Result), lengthOfGitProviders)
	})

}
