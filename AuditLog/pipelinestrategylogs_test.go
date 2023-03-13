package AuditLog

import (
	"automation-suite/PipelineConfigRouter"
	"automation-suite/PipelineConfigRouter/RequestDTOs"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"strings"
	"time"
)

func (suite *AuditLogsTestSuite) TestClassToCheckPipelineStrategyLogs() {

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

	workflowResponse := PipelineConfigRouter.HitCreateWorkflowApiWithFullPayload(createAppApiResponse.Id, suite.authToken).Result
	preStageScript, _ := Base.GetByteArrayOfGivenJsonFile("../testdata/PipeLineConfigRouter/preStageScript.txt")
	postStageScript, _ := Base.GetByteArrayOfGivenJsonFile("../testdata/PipeLineConfigRouter/postStageScript.txt")

	suite.Run("ABCD", func() {

		payload := PipelineConfigRouter.GetRequestPayloadForSaveCdPipelineApi(createAppApiResponse.Id, workflowResponse.AppWorkflowId, 1, workflowResponse.CiPipelines[0].Id, workflowResponse.CiPipelines[0].ParentCiPipeline, Automatic, string(preStageScript), string(postStageScript), Manual)
		bytePayload, _ := json.Marshal(payload)
		savePipelineResponse := PipelineConfigRouter.HitSaveCdPipelineApi(bytePayload, suite.authToken)
		time.Sleep(2 * time.Second)

		query1 := "SELECT * FROM pipeline_strategy_history ORDER BY Id DESC LIMIT 1"
		row1 := Base.GetData(query1)

		var firstRow1 PipelineStrategyTableHistory

		for row1.Next() {
			row := row1.Scan(&firstRow1.Id, &firstRow1.PipelineId, &firstRow1.Strategy, &firstRow1.Config, &firstRow1.Default,
				&firstRow1.Deployed, &firstRow1.DeployedOn, &firstRow1.DeployedBy,
				&firstRow1.CreatedOn, &firstRow1.CreatedBy, &firstRow1.UpdatedOn, &firstRow1.UpdatedBy, &firstRow1.PipelineTriggerType)
			if row != nil {
				log.Fatal(row)
			}
		}

		err := row1.Err()
		if err != nil {
			log.Fatal(err)
		}

		var pipelinesArray []RequestDTOs.Pipeline

		pipelinesArray = payload.Pipelines
		pipeline := pipelinesArray[0]
		tempId := savePipelineResponse.Result.Pipelines[0]

		assert.Equal(suite.T(), pipeline.TriggerType, firstRow1.PipelineTriggerType, "These should match , values either manual or automatic")
		assert.Equal(suite.T(), firstRow1.PipelineId, tempId.Id, "These should match , pipelineId would be same")
		assert.True(suite.T(), strings.Contains(firstRow1.Strategy, "rolling"), "The default strategy is rolling strategy")

		deletePipelinePayload := PipelineConfigRouter.GetPayloadForDeleteCdPipeline(createAppApiResponse.Id, savePipelineResponse.Result.Pipelines[0].Id)
		deletePipelineByteCode, _ := json.Marshal(deletePipelinePayload)
		PipelineConfigRouter.HitForceDeleteCdPipelineApi(deletePipelineByteCode, suite.authToken)
	})

	log.Println("=== Deleting the CI pipeline ===")
	PipelineConfigRouter.DeleteCiPipeline(createAppApiResponse.Id, workflowResponse.CiPipelines[0].Id, suite.authToken)
	log.Println("=== CI Workflow ===")
	PipelineConfigRouter.HitDeleteWorkflowApi(createAppApiResponse.Id, workflowResponse.AppWorkflowId, suite.authToken)
	log.Println("=== Deleting the Test data created after verification ===")
	Base.DeleteApp(createAppApiResponse.Id, createAppApiResponse.AppName, createAppApiResponse.TeamId, createAppApiResponse.TemplateId, suite.authToken)
}
