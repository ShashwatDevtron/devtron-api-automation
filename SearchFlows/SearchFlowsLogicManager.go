package SearchFlows

import (
	"automation-suite/EnvironmentRouter"
	Request "automation-suite/EnvironmentRouter/RequestDTOs"
	EnvironmentRouterResponseDTOs "automation-suite/EnvironmentRouter/ResponseDTOs"
	"automation-suite/PipelineConfigRouter"
	"automation-suite/TeamRouter"
	TeamRouterResponseDTOs "automation-suite/TeamRouter/ResponseDTOs"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/suite"
)

func CreateProject(payload []byte, authToken string) TeamRouterResponseDTOs.SaveTeamResponseDTO {
	response := TeamRouter.HitSaveTeamApi(payload, authToken)
	return response
}

func DeleteProject(payload []byte, authToken string) TeamRouterResponseDTOs.DeleteTeamResponseDto {
	response := TeamRouter.HitDeleteTeamApi(payload, authToken)
	return response
}

func CreateEnv(payload []byte, authToken string) EnvironmentRouterResponseDTOs.CreateEnvironmentResponseDTO {
	response := EnvironmentRouter.HitCreateEnvApi(payload, authToken)
	return response
}

func DeleteEnv(payload []byte, authToken string) EnvironmentRouterResponseDTOs.DeleteEnvResponseDTO {
	response := EnvironmentRouter.HitDeleteEnvApi(payload, authToken)
	return response
}

func CreateDevtronApp(appName string, authToken string, teamId int) PipelineConfigRouter.CreateAppResponseDto {
	//appName := "app" + strings.ToLower(Base.GetRandomStringOfGivenLength(10))
	createAppRequestDto := PipelineConfigRouter.GetAppRequestDto(appName, teamId, 0)
	byteValueOfCreateApp, _ := json.Marshal(createAppRequestDto)
	response := PipelineConfigRouter.HitCreateAppApi(byteValueOfCreateApp, appName, teamId, 0, authToken)
	return response
}
func DeleteDevtronApp(appId int, appName string, teamId int, TemplateId int, authToken string) PipelineConfigRouter.DeleteResponseDto {
	byteValueOfDeleteApp := PipelineConfigRouter.GetPayLoadForDeleteAppAPI(appId, appName, teamId, TemplateId)
	response := PipelineConfigRouter.HitDeleteAppApi(byteValueOfDeleteApp, appId, authToken)
	return response
}

func GetSaveEnvRequestDto() Request.CreateEnvironmentRequestDTO {
	var saveEnvRequestDto Request.CreateEnvironmentRequestDTO
	EnvName := Base.GetRandomStringOfGivenLength(10)
	saveEnvRequestDto.Environment = EnvName
	saveEnvRequestDto.Active = true
	namespace := Base.GetRandomStringOfGivenLengthOfLowerCaseAndNumber(10)
	saveEnvRequestDto.Namespace = namespace
	saveEnvRequestDto.ClusterId = 1

	return saveEnvRequestDto
}

type SearchFlowTestSuite struct {
	suite.Suite
	authToken string
}

func (suite *SearchFlowTestSuite) SetupSuite() {
	suite.authToken = Base.GetAuthToken()
}
