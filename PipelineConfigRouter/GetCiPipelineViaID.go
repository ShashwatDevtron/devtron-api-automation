package PipelineConfigRouter

import (
	"automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"strings"
)

func (suite *PipelineConfigSuite) TestClass3GetCiPipeline() {
	config, _ := GetEnvironmentConfigPipelineConfigRouter()
	createAppApiResponse := suite.createAppResponseDto.Result
	materialId := suite.createAppMaterialResponseDto.Result.Material[0].Id
	appNameForNewCreation := "app" + strings.ToLower(testUtils.GetRandomStringOfGivenLength(6))

	suite.Run("A=8=GetCiPipelineWithoutCreatingIt", func() {
		log.Println("=== Here We are creating a new App ===")
		createAppRequestDto := GetAppRequestDto(appNameForNewCreation, 1, 0)
		byteValueOfCreateApp, _ := json.Marshal(createAppRequestDto)
		createAppResponseDto := HitCreateAppApi(byteValueOfCreateApp, appNameForNewCreation, 1, 0, suite.authToken)

		log.Println("=== Hitting the GetCiPipelineViaValidID API ====")
		getCiPipelineResponse := HitGetCiPipelineViaId(strconv.Itoa(createAppResponseDto.Result.Id), suite.authToken)
		assert.Equal(suite.T(), getCiPipelineResponse.Errors[0].UserMessage, "no ci pipeline exists")
	})

	suite.Run("A=9=GetCiPipelineViaValidAppID", func() {
		log.Println("=== getting Test Data for Hitting the SaveAppCiPipeline API ====")
		appName := createAppApiResponse.AppName
		requestPayloadForSaveAppCiPipeline := getRequestPayloadForSaveAppCiPipeline(createAppApiResponse.Id, config.DockerRegistry, config.DockerRegistry+"/"+appName, config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, materialId)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		log.Println("=== Hitting the SaveAppCiPipeline API ====")
		HitSaveAppCiPipeline(byteValueOfSaveAppCiPipeline, suite.authToken)
		log.Println("=== Hitting the GetCiPipelineViaValidID API ====")
		getCiPipelineResponse := HitGetCiPipelineViaId(strconv.Itoa(createAppApiResponse.Id), suite.authToken)
		assert.Equal(suite.T(), getCiPipelineResponse.Result.AppName, createAppApiResponse.AppName)
		assert.Equal(suite.T(), getCiPipelineResponse.Result.DockerBuildConfig.GitMaterialId, materialId)
	})

	suite.Run("A=10=GetCiPipelineViaInValidAppID", func() {
		log.Println("=== getting Test Data for Hitting the GetCiPipelineViaInValidAppID API ====")
		invalidAppId := strconv.Itoa(testUtils.GetRandomNumberOf9Digit())
		log.Println("=== Hitting the GetCiPipelineViaInValidAppID API ====")
		getCiPipelineResponse := HitGetCiPipelineViaId(invalidAppId, suite.authToken)
		assert.Equal(suite.T(), getCiPipelineResponse.Errors[0].UserMessage, "pg: no rows in result set")
	})
	// <tear-down code>
}
