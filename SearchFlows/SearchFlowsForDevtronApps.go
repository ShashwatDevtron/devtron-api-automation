package SearchFlows

import (
	"automation-suite/AppListingRouter"
	"automation-suite/AppListingRouter/ResponseDTOs"
	"automation-suite/SearchFlows/RequestDTOs"
	"automation-suite/TeamRouter"
	"automation-suite/UserRouter"
	"encoding/json"
	"log"
	"strings"

	"github.com/stretchr/testify/assert"
)

func getStatusCheck(expectedCode, actualCode int) bool {
	if expectedCode == 403 && actualCode == 200 {
		return false
	} else if expectedCode == 200 && (actualCode == 403 || actualCode == 401) {
		return false
	}
	return true
}

func implContains(sl []ResponseDTOs.AppContainers, name string) bool {
	for _, value := range sl {
		if value.AppName == name {
			return true
		}
	}
	return false
}

func (suite *SearchFlowTestSuite) TestRbacFlowsForDevtronApps() {

	suite.Run("A=0=SearchAppWithCreatingNewApplication", func() {
		// Creating Project with Super Admin
		var devtronDeletion RequestDTOs.SearchDevtronDeletion
		saveTeamRequestDto := TeamRouter.GetSaveTeamRequestDto()
		saveTeamRequestDto.Name = UserRouter.PROJECT
		byteValueOfStruct, _ := json.Marshal(saveTeamRequestDto)

		responseOfCreateProject := CreateProject(byteValueOfStruct, suite.authToken)
		assert.Equal(suite.T(), 200, responseOfCreateProject.Code)
		assert.Equal(suite.T(), UserRouter.PROJECT, responseOfCreateProject.Result.Name)
		devtronDeletion.ProjectPayload, _ = json.Marshal(responseOfCreateProject.Result)

		// Creating environment with SuperAdmin
		environments := strings.Split(UserRouter.ENV, ",")
		saveEnvRequestDto := GetSaveEnvRequestDto()
		saveEnvRequestDto.Environment = environments[0]
		saveEnvRequestDto.EnvironmentIdentifier = "default_cluster__" + saveEnvRequestDto.Namespace
		byteValueOfStruct, _ = json.Marshal(saveEnvRequestDto)
		responseOfCreateEnvironment := CreateEnv(byteValueOfStruct, suite.authToken)
		assert.Equal(suite.T(), 200, responseOfCreateEnvironment.Code)
		assert.Equal(suite.T(), environments[0], responseOfCreateEnvironment.Result.Environment)

		devtronDeletion.EnvPayLoad, _ = json.Marshal(responseOfCreateEnvironment.Result)

		//Application With SuperAdmin
		applications := strings.Split(UserRouter.APP, ",")
		responseOfCreateDevtronApp := CreateDevtronApp(applications[0], suite.authToken, responseOfCreateProject.Result.Id)
		assert.Equal(suite.T(), 200, responseOfCreateDevtronApp.Code)
		assert.Equal(suite.T(), applications[0], responseOfCreateDevtronApp.Result.AppName)
		assert.Equal(suite.T(), responseOfCreateProject.Result.Id, responseOfCreateDevtronApp.Result.TeamId)
		devtronDeletion.DevtronPayload = responseOfCreateDevtronApp

		Envs := []int{}
		Teams := []int{responseOfCreateProject.Result.Id}
		Namespaces := []string{}
		AppStatuses := []string{}
		requestDTOForApiFetchAppsByEnvironment := AppListingRouter.GetPayloadForApiFetchAppsByEnvironment(Envs, Teams, Namespaces, "", AppStatuses, "ASC", 0, 0, 10)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(requestDTOForApiFetchAppsByEnvironment)

		log.Println("Test Case for SuperAdmin ===>", suite.authToken)
		allAppsByEnvironment := AppListingRouter.HitApiFetchAppsByEnvironment(bytePayloadForTriggerCiPipeline, suite.authToken)

		assert.Equal(suite.T(), true, getStatusCheck(200, allAppsByEnvironment.Code))
		//to check if the app is created
		testAppCount := len(strings.Split(UserRouter.APP, ","))
		assert.Equal(suite.T(), testAppCount, allAppsByEnvironment.Result.AppCount)
		assert.Equal(suite.T(), UserRouter.APP, allAppsByEnvironment.Result.AppContainers[0].AppName)

		// delete the app
		responseOfDeleteProject := DeleteDevtronApp(devtronDeletion.DevtronPayload.Result.Id, devtronDeletion.DevtronPayload.Result.AppName, devtronDeletion.DevtronPayload.Result.TeamId, devtronDeletion.DevtronPayload.Result.TemplateId, suite.authToken)
		assert.Equal(suite.T(), true, getStatusCheck(200, responseOfDeleteProject.Code))

		//to check if the app is deleted
		allAppsByEnvironment = AppListingRouter.HitApiFetchAppsByEnvironment(bytePayloadForTriggerCiPipeline, suite.authToken)
		assert.Equal(suite.T(), testAppCount-1, allAppsByEnvironment.Result.AppCount)
		assert.Equal(suite.T(), false, implContains(allAppsByEnvironment.Result.AppContainers, devtronDeletion.DevtronPayload.Result.AppName))

		DeleteProject(devtronDeletion.ProjectPayload, suite.authToken)
	})

}
