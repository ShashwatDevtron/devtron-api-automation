package GitProviderRouter

import (
	"automation-suite/GitProviderRouter/ResponseDTOs"
	"automation-suite/PipelineConfigRouter"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
)

func (suite *GitProRouterTestSuite) TestClassA6GetApp() {

	var byteValueOfSaveGitProvider []byte
	var saveGitProviderResponseDto ResponseDTOs.SaveGitProviderResponseDto
	var getGitProviderResponseLength int

	suite.Run("getting payload for get Git Provider api", func() {
		log.Println("Hitting The get git provider API")
		getGitProviderResponse := HitGetGitProviderApi(suite.authToken)
		getGitProviderResponseLength = len(getGitProviderResponse.Result)

		log.Println("SaveGitProviderWithValidPayload")

		saveGitProviderRequestDto := GetGitProviderRequestDto(1)
		byteValueOfSaveGitProvider, _ = json.Marshal(saveGitProviderRequestDto)

		log.Println("Hitting The post git provider API")
		saveGitProviderResponseDto = HitSaveGitProviderApi(byteValueOfSaveGitProvider, suite.authToken)

		log.Println("Validating the Response of the save git provider API...")
		assert.Equal(suite.T(), saveGitProviderRequestDto.Name, saveGitProviderResponseDto.Result.Name)

		log.Println("Hitting The get git provider API")
		getGitProviderResponses := HitGetGitProviderApi(suite.authToken)
		assert.Equal(suite.T(), len(getGitProviderResponses.Result), getGitProviderResponseLength+1)

		//saveGitProviderRequestPayload := GetGitProviderRequestDto(1)
		//byteValueOfGitProviderPayload, _ := json.Marshal(saveGitProviderRequestPayload)

		//log.Println("Hitting The save git provider Api First time")
		//saveGitProviderResponse := HitSaveGitProviderApi(byteValueOfGitProviderPayload, suite.authToken)

		//log.Println("getting payload for git provider API")
		//byteValueOfDeleteGitProvider := GetPayLoadForDeleteGitProviderAPI(saveGitProviderResponse.Result.Id, saveGitProviderResponse.Result.GitHostId, saveGitProviderResponse.Result.Url, saveGitProviderResponse.Result.AuthMode, saveGitProviderResponse.Result.Name)
		//log.Println("Hitting the Delete team API for Removing the data created via automation")
		//HitDeleteGitProviderApi(byteValueOfDeleteGitProvider, suite.authToken)

		//log.Println("Hitting The get git provider by id API")
		//getGitProviderByIdResponse := HitGetGitProviderByIdApi(saveGitProviderResponseDto.Result.Id, suite.authToken)
		//assert.Equal(suite.T(), getGitProviderByIdResponse.Result.Name, saveGitProviderResponseDto.Result.Name)

		envConf := Base.ReadBaseEnvConfig()
		config := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)

		log.Println("=== Here we are creating a App ===")
		createAppApiResponse := Base.CreateApp(suite.authToken).Result

		log.Println("=== Here we are creating App Material ===")
		createAppMaterialResponse := PipelineConfigRouter.HitCreateAppMaterialApi(nil, createAppApiResponse.Id, saveGitProviderResponseDto.Result.Id, false, suite.authToken)

		log.Println("=== Saving docker build config ===")
		requestPayloadForSaveAppCiPipeline := PipelineConfigRouter.GetRequestPayloadForSaveAppCiPipeline(createAppApiResponse.Id, config.DockerRegistry, config.DockerRegistry+"/test", config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, createAppMaterialResponse.Result.Material[0].Id)
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

		workflowResponse := PipelineConfigRouter.HitCreateWorkflowApiWithFullPayload(createAppApiResponse.Id, suite.authToken).Result
		log.Println("=== Hitting GetCiPipelineMaterial API ===")
		pipelineMaterial := PipelineConfigRouter.HitGetCiPipelineMaterial(workflowResponse.CiPipelines[0].Id, suite.authToken)
		payloadForTriggerCiPipeline := PipelineConfigRouter.CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.CiPipelines[0].Id, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		log.Println("=== Hitting TriggerCiPipeline API ===")
		_ = PipelineConfigRouter.HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)

		log.Println("getting payload for Update git provider API")
		byteValueOfUpdateGitProvider := GetPayLoadForUpdateGitProviderAPI(saveGitProviderResponseDto.Result.Id, saveGitProviderResponseDto.Result.GitHostId)
		log.Println("Hitting the Update Git API for Updating the data created via automation")
		HitUpdateGitProviderApi(byteValueOfUpdateGitProvider, suite.authToken)

		log.Println("getting payload for Delete git provider API")
		byteValueOfDeleteDockerRegistry := GetPayLoadForDeleteGitProviderAPI(saveGitProviderResponseDto.Result.Id, saveGitProviderResponseDto.Result.GitHostId, saveGitProviderResponseDto.Result.Url, saveGitProviderResponseDto.Result.AuthMode, saveGitProviderResponseDto.Result.Name)
		log.Println("Hitting the Delete Git API for Removing the data created via automation")
		HitDeleteGitProviderApi(byteValueOfDeleteDockerRegistry, suite.authToken)

		log.Println("Hitting The get git provider API")
		getGitProviderResponseses := HitGetGitProviderApi(suite.authToken)
		assert.Equal(suite.T(), len(getGitProviderResponseses.Result), getGitProviderResponseLength)
	})

	suite.Run("A=2=SaveGitProviderWithExistingName", func() {
		saveGitProviderRequestPayload := GetGitProviderRequestDto(1)
		byteValueOfGitProviderPayload, _ := json.Marshal(saveGitProviderRequestPayload)

		log.Println("Hitting The save git provider Api second time with existing registry name")
		finalApiResponse := HitSaveGitProviderApi(byteValueOfGitProviderPayload, suite.authToken)

		log.Println("Validating the Response of the save git provider  API...")
		assert.Equal(suite.T(), finalApiResponse.Status, "Internal Server Error")
		assert.Equal(suite.T(), finalApiResponse.Status, "Internal Server Error")

	})

}
