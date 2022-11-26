package AppStoreDiscoverRouter

import (
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/sjson"
	"log"
	"sigs.k8s.io/yaml"
	"strconv"
	"time"
)

func (suite *AppStoreDiscoverTestSuite) TestDiscoverPreviouslyInstalledHelmAppsViaRepoId() {
	var valuesOverrideInterface interface{}
	var AppStoreId int
	log.Println("=== Getting apache chart repo via DiscoverApp API ===")
	queryParams := map[string]string{"appStoreName": "apache"}
	PollForGettingHelmAppData(queryParams, suite.authToken)
	ActiveDiscoveredApps := HitDiscoverAppApi(queryParams, suite.authToken)
	var requiredReferenceId int
	for _, DiscoveredApp := range ActiveDiscoveredApps.Result {
		if DiscoveredApp.ChartName == "bitnami" {
			requiredReferenceId = DiscoveredApp.AppStoreApplicationVersionId
			AppStoreId = DiscoveredApp.Id
			break
		}
	}
	log.Println("=== Getting Template values for apache chart===")
	queryParamsOfApi := map[string]string{"referenceId": strconv.Itoa(requiredReferenceId), "kind": "DEFAULT"}
	referenceTemplate := HitGetTemplateValuesViaReferenceIdApi(queryParamsOfApi, suite.authToken)
	valuesOverrideYaml := referenceTemplate.Result.Values
	if err := yaml.Unmarshal([]byte(valuesOverrideYaml), &valuesOverrideInterface); err != nil {
		panic(err)
	}
	Base.ConvertYamlIntoJson(valuesOverrideInterface)
	valuesOverrideJson, _ := json.Marshal(valuesOverrideInterface)
	jsonOfSaveDeploymentTemp := string(valuesOverrideJson)
	jsonWithTypeAsClusterIP, _ := sjson.Set(jsonOfSaveDeploymentTemp, "service.type", "ClusterIP")
	updatedValuesOverrideJson := []byte(jsonWithTypeAsClusterIP)
	log.Println("=== converting Json into YAML for Values Override in Install API===")
	updatedValuesOverrideYaml, _ := yaml.JSONToYAML(updatedValuesOverrideJson)
	var installedAppId int
	installAppRequestDTO := GetRequestDtoForInstallApp(requiredReferenceId, requiredReferenceId, valuesOverrideInterface, string(updatedValuesOverrideYaml))
	byteValueOfInstallAppRequestPayload, _ := json.Marshal(installAppRequestDTO)
	jsonOfSaveDeploymentTemp1 := string(byteValueOfInstallAppRequestPayload)
	jsonWithTypeAsClusterIP1, _ := sjson.Set(jsonOfSaveDeploymentTemp1, "valuesOverride.service.type", "ClusterIP")
	updatedByteValueOfInstallAppRequestPayload := []byte(jsonWithTypeAsClusterIP1)
	responseAfterInstallingApp := HitInstallAppApi(string(updatedByteValueOfInstallAppRequestPayload), suite.authToken)
	installedAppId = responseAfterInstallingApp.Result.InstalledAppId
	time.Sleep(10 * time.Second)
	suite.Run("A=1=GetInstalledAppsByAppStoreId", func() {
		log.Println("Hitting the GetDeploymentOfInstalledApp API with valid payload")
		deploymentOfInstalledApp := GetInstalledAppsByAppStoreId(strconv.Itoa(AppStoreId), suite.authToken)
		assert.NotNil(suite.T(), deploymentOfInstalledApp.Result[len(deploymentOfInstalledApp.Result)-1].InstalledAppVersionId)
		assert.Equal(suite.T(), installedAppId, deploymentOfInstalledApp.Result[len(deploymentOfInstalledApp.Result)-1].InstalledAppId)
	})
	log.Println("Removing the data created via API")
	respOfDeleteInstallAppApi := HitDeleteInstalledAppApi(strconv.Itoa(installedAppId), suite.authToken)
	assert.Equal(suite.T(), installedAppId, respOfDeleteInstallAppApi.Result.InstalledAppId)
}