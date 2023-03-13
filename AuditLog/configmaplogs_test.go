package AuditLog

import (
	"automation-suite/HelperRouter"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
)

func (suite *AuditLogsTestSuite) TestClassToCheckConfigMapLogs() {
	createAppApiResponse := Base.CreateApp(suite.authToken).Result

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
