package AuditLog

import (
	"automation-suite/HelperRouter"
	"automation-suite/PipelineConfigRouter"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"strings"
)

func (suite *AuditLogsTestSuite) TestClassToCheckConfigMapLogs() {
	envConf := Base.ReadBaseEnvConfig()
	config := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)

	log.Println("=== Creating a App ===")
	createAppApiResponse := Base.CreateApp(suite.authToken).Result

	log.Println("=== Creating App Material ===")
	createAppMaterialRequestDto := PipelineConfigRouter.GetAppMaterialRequestDto(createAppApiResponse.Id, 1, false)
	appMaterialByteValue, _ := json.Marshal(createAppMaterialRequestDto)
	createAppMaterialResponse := PipelineConfigRouter.HitCreateAppMaterialApi(appMaterialByteValue, createAppApiResponse.Id, 1, false, suite.authToken)

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

	suite.Run("ABC", func() {

		configName := strings.ToLower(Base.GetRandomStringOfGivenLength(6))
		requestPayloadForSecret := HelperRouter.GetRequestPayloadForSecretOrConfig(0, configName, createAppApiResponse.Id, "environment", "kubernetes", false, false, false, true)
		byteRequestPayloadForSecret, _ := json.Marshal(requestPayloadForSecret)
		HelperRouter.HitSaveGlobalConfigMap(byteRequestPayloadForSecret, suite.authToken)

		query1 := "SELECT * FROM config_map_history ORDER BY Id DESC LIMIT 1"
		row1 := Base.GetData(query1)

		var firstRow1 ConfigMapTableHistory

		for row1.Next() {
			row := row1.Scan(&firstRow1.Id, &firstRow1.PipelineId, &firstRow1.AppId, &firstRow1.DataType, &firstRow1.Data,
				&firstRow1.Deployed, &firstRow1.DeployedOn, &firstRow1.DeployedBy, &firstRow1.CreatedOn,
				&firstRow1.CreatedBy, &firstRow1.UpdatedOn, &firstRow1.UpdatedBy)
			if row != nil {
				log.Fatal(row)
			}
		}

		err := row1.Err()
		if err != nil {
			log.Fatal(err)
		}

		//assert.Equal(suite.T(), requestPayloadForSecret.Id, firstRow1.Id, "The two arguments should be the same.")
		assert.Equal(suite.T(), requestPayloadForSecret.AppId, firstRow1.AppId, "The two arguments should be the same.")
		assert.True(suite.T(), strings.Contains(firstRow1.Data, configName), "The Appname from database should contain configname")
	})

	log.Println("=== Here We are Deleting the test data created for Automation ===")
	Base.DeleteApp(createAppApiResponse.Id, createAppApiResponse.AppName, createAppApiResponse.TeamId, createAppApiResponse.TemplateId, suite.authToken)
}
