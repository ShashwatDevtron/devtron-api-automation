package GitopsFlows

import (
	"automation-suite/GitopsConfigRouter"
	GitopsConfigRouterResponseDTOs "automation-suite/GitopsConfigRouter/ResponseDTOs"
	Base "automation-suite/testUtils"
	"github.com/stretchr/testify/suite"
)

func SaveGitopsConfig(payload []byte, provider string, username string, host string, token string, githuborgid string, authToken string) GitopsConfigRouterResponseDTOs.CreateGitopsConfigResponseDto {
	return GitopsConfigRouter.HitCreateGitopsConfigApi(payload, provider, username, host, token, githuborgid, authToken)
}
func ValidateGitops(payload []byte, provider string, username string, host string, token string, githuborgid string, authToken string) GitopsConfigRouterResponseDTOs.CreateGitopsConfigResponseDto {
	return GitopsConfigRouter.HitValidateGitopsConfigApi(payload, provider, username, host, token, githuborgid, authToken)
}
func FetchAllGitopsConfig(authToken string) GitopsConfigRouterResponseDTOs.FetchAllGitopsConfigResponseDto {
	return GitopsConfigRouter.HitFetchAllGitopsConfigApi(authToken)
}
func IsGitopsConfigured(authToken string) GitopsConfigRouterResponseDTOs.CheckGitopsExistsResponse {
	return GitopsConfigRouter.HitGitopsConfigured(authToken)
}

type GitopsFlowsTestSuite struct {
	suite.Suite
	authToken string
}

func (suite *GitopsFlowsTestSuite) SetupSuite() {
	suite.authToken = Base.GetAuthToken()
}
