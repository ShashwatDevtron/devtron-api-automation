package PipelineConfigRouter

import (
	"automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
)

func (suite *PipelinesConfigRouterTestSuite) TestClass2SaveAppCiPipeline() {
	config, _ := GetEnvironmentConfigPipelineConfigRouter()
	createAppApiResponse := suite.createAppResponseDto.Result
	createAppMaterialResponse := suite.createAppMaterialResponseDto.Result

	suite.Run("A=1=SaveAppCiPipelineWithValidPayload", func() {
		requestPayloadForSaveAppCiPipeline := getRequestPayloadForSaveAppCiPipeline(createAppApiResponse.Id, config.DockerRegistry, config.DockerRegistry+"/test", config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, createAppMaterialResponse.Material[0].Id)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		log.Println("=== Hitting the SaveAppCiPipeline API ====")
		saveAppCiPipelineResponse := HitSaveAppCiPipeline(byteValueOfSaveAppCiPipeline, suite.authToken)
		assert.Equal(suite.T(), saveAppCiPipelineResponse.Result.AppName, createAppApiResponse.AppName)
	})

	suite.Run("A=2=SaveAppCiPipelineWithInValidAppId", func() {
		appId := testUtils.GetRandomNumberOf9Digit()
		requestPayloadForSaveAppCiPipeline := getRequestPayloadForSaveAppCiPipeline(appId, config.DockerRegistry, config.DockerRegistry+"/test", config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, createAppMaterialResponse.Material[0].Id)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		log.Println("=== Hitting the SaveAppCiPipeline API with Invalid AppId ====")
		saveAppCiPipelineResponse := HitSaveAppCiPipeline(byteValueOfSaveAppCiPipeline, suite.authToken)
		assert.Equal(suite.T(), saveAppCiPipelineResponse.Errors[0].UserMessage, "pg: no rows in result set")
	})

	suite.Run("A=3=SaveAppCiPipelineWithInValidMaterialId", func() {
		invalidMaterialId := testUtils.GetRandomNumberOf9Digit()
		requestPayloadForSaveAppCiPipeline := getRequestPayloadForSaveAppCiPipeline(createAppApiResponse.Id, config.DockerRegistry, config.DockerRegistry+"/test", config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, invalidMaterialId)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		log.Println("=== Hitting the SaveAppCiPipeline API with Invalid Material Id ====")
		saveAppCiPipelineResponse := HitSaveAppCiPipeline(byteValueOfSaveAppCiPipeline, suite.authToken)
		assert.Equal(suite.T(), saveAppCiPipelineResponse.Errors[0].UserMessage, "ERROR #23505 duplicate key value violates unique constraint \"ci_template_app_id_key\"")
	})

	suite.Run("A=4=SaveAppCiPipelineWithInValidDockerfileRepository", func() {
		requestPayloadForSaveAppCiPipeline := getRequestPayloadForSaveAppCiPipeline(createAppApiResponse.Id, config.DockerRegistry+"invalid", config.DockerRegistry+"/test", config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, createAppMaterialResponse.Material[0].Id)
		byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
		log.Println("=== Hitting the SaveAppCiPipeline API with Invalid Docker file Repository ====")
		saveAppCiPipelineResponse := HitSaveAppCiPipeline(byteValueOfSaveAppCiPipeline, suite.authToken)
		assert.Equal(suite.T(), saveAppCiPipelineResponse.Errors[0].UserMessage, "pg: no rows in result set")
	})
}
