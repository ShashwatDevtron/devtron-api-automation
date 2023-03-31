package GitProviderFullFlow

import (
	"automation-suite/GitProviderRouter"
	"automation-suite/GitProviderRouter/RequestDTOs"
	"automation-suite/GitProviderRouter/ResponseDTOs"
	"automation-suite/PipelineConfigRouter"
	CiPipelineRequestDTO "automation-suite/PipelineConfigRouter/RequestDTOs"
	CiPipelineResponseDTO "automation-suite/PipelineConfigRouter/ResponseDTOs"
	Base "automation-suite/testUtils"
	"automation-suite/testdata/testUtils"
	"fmt"
	"github.com/stretchr/testify/suite"
)

func HitGetListOfGitProviders(authToken string) ResponseDTOs.GetGitProviderResponseDto {
	return GitProviderRouter.HitGetGitProviderApi(authToken)
}

func GetGitProviderRequestObjectDto(GitRegHostId int, authMode string) RequestDTOs.SaveGitProviderRequestDTO {
	return GitProviderRouter.GetGitProviderRequestDto(GitRegHostId, authMode)
}

func HitSaveOneGitProviderApi(payloadOfApi []byte, authToken string) ResponseDTOs.SaveGitProviderResponseDto {
	return GitProviderRouter.HitSaveGitProviderApi(payloadOfApi, authToken)
}

func CreateAppWithRandomData(authToken string) testUtils.CreateAppResponseDto {
	return testUtils.CreateApp(authToken)
}

func HitCreateAppMaterialWithGitProviderIdApi(payload []byte, appId int, gitProviderId int, fetchSubmodules bool, authToken string) PipelineConfigRouter.CreateAppMaterialResponseDto {
	return PipelineConfigRouter.HitCreateAppMaterialApi(payload, appId, gitProviderId, fetchSubmodules, authToken)
}

func GetRequestPayloadForSavingAppCiPipeline(AppId int, dockerRegistry string, dockerRepository string, dockerfilePath string, dockerfileRepository string, dockerfileRelativePath string, gitMaterialId int) CiPipelineRequestDTO.SaveAppCiPipelineRequestDTO {
	return PipelineConfigRouter.GetRequestPayloadForSaveAppCiPipeline(AppId, dockerRegistry, dockerRepository, dockerfilePath, dockerfileRepository, dockerfileRelativePath, gitMaterialId)
}

func HitSaveCiPipelineForApp(payload []byte, authToken string) PipelineConfigRouter.SaveAppCiPipelineResponseDTO {
	return PipelineConfigRouter.HitSaveAppCiPipeline(payload, authToken)
}

func HitGetChartRefViaAppId(appId string, authToken string) PipelineConfigRouter.GetChartReferenceResponseDTO {
	return PipelineConfigRouter.HitGetChartReferenceViaAppId(appId, authToken)
}

func HitGetTemplateViaAppIdAndChartRef(appId string, chartRefId string, authToken string) PipelineConfigRouter.GetAppTemplateResponseDto {
	return PipelineConfigRouter.HitGetTemplateViaAppIdAndChartRefId(appId, chartRefId, authToken)
}

func GetRequestPayloadForSavingDeploymentTemplate(AppId int, chartRefId int, defaultOverride PipelineConfigRouter.DefaultAppOverride) PipelineConfigRouter.SaveDeploymentTemplateRequestDTO {
	return PipelineConfigRouter.GetRequestPayloadForSaveDeploymentTemplate(AppId, chartRefId, defaultOverride)
}

func HitSavingDeploymentTemplateApi(payload []byte, authToken string) PipelineConfigRouter.SaveDeploymentTemplateResponseDTO {
	return PipelineConfigRouter.HitSaveDeploymentTemplateApi(payload, authToken)
}

func HitCreatingWorkflowApiWithFullPayload(appId int, authToken string) CiPipelineResponseDTO.CreateWorkflowResponseDto {
	return PipelineConfigRouter.HitCreateWorkflowApiWithFullPayload(appId, authToken)
}

func HitGettingCiPipelineMaterial(ciPipelineId int, authToken string) CiPipelineResponseDTO.GetCiPipelineMaterialResponseDTO {
	return PipelineConfigRouter.HitGetCiPipelineMaterial(ciPipelineId, authToken)
}
func CreatePayloadForTriggeringCiPipeline(commit string, PipelineId int, ciPipelineMaterialId int, invalidateCache bool) CiPipelineRequestDTO.TriggerCiPipelineRequestDTO {
	return PipelineConfigRouter.CreatePayloadForTriggerCiPipeline(commit, PipelineId, ciPipelineMaterialId, invalidateCache)
}

func HittingTriggerCiPipelineApi(payload []byte, authToken string) CiPipelineResponseDTO.TriggerCiPipelineResponseDTO {
	return PipelineConfigRouter.HitTriggerCiPipelineApi(payload, authToken)
}

func HitGetOneGitProviderByIdApi(appId int, authToken string) ResponseDTOs.GetGitProviderResponseById {
	return GitProviderRouter.HitGetGitProviderByIdApi(appId, authToken)
}

func GetPayLoadForUpdateOneGitProviderAPI(id int, gitHostId int) []byte {
	return GitProviderRouter.GetPayLoadForUpdateGitProviderAPI(id, gitHostId)
}

func HitUpdateOneGitProviderApi(byteValueOfStruct []byte, authToken string) ResponseDTOs.SaveGitProviderResponseDto {
	return GitProviderRouter.HitUpdateGitProviderApi(byteValueOfStruct, authToken)
}

func GetPayLoadForDeleteOneGitProviderAPI(id int, gitHostId int, url string, authMode string, name string) []byte {
	return GitProviderRouter.GetPayLoadForDeleteGitProviderAPI(id, gitHostId, url, authMode, name)
}

func HitDeleteOneGitProviderApi(byteValueOfStruct []byte, authToken string) ResponseDTOs.DeleteGitProviderResponse {
	return GitProviderRouter.HitDeleteGitProviderApi(byteValueOfStruct, authToken)
}

type GitProRouterTestSuite struct {
	suite.Suite
	authToken string
}

func (suite *GitProRouterTestSuite) SetupSuite() {
	suite.authToken = Base.GetAuthToken()
}

func (suite *GitProRouterTestSuite) checkForCiLogs(pipelineId string, ciWorkflowId string, LogLineNumber int) {
	ciLogsDownloadUrlFormat := PipelineConfigRouter.GetCiPipelineBaseUrl + "/%s/workflow/%s/logs"
	ciLogsDownloadUrl := fmt.Sprintf(ciLogsDownloadUrlFormat, pipelineId, ciWorkflowId)
	testUtils.ReadEventStreamsForSpecificApiAndVerifyResult(ciLogsDownloadUrl, suite.authToken, suite.T(), LogLineNumber, "git cloning", true)
}
