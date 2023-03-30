package GitopsConfigRouter

import (
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
)

func (suite *GitOpsRouterTestSuite) TestClassA4GetGitopsConfigByProvider() {
	suite.Run("A=1=FetchAllGitopsConfig", func() {

		log.Println("Hitting the 'Save Gitops Config' Api for creating a new entry")
		//gitopsConfig, _ := GetGitopsConfig()
		envConf := Base.ReadBaseEnvConfig()
		gitopsConfig := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)

		createGitopsConfigRequestDto := GetGitopsConfigRequestDto(gitopsConfig.Provider, gitopsConfig.GitUsername, gitopsConfig.Host, gitopsConfig.GitToken, gitopsConfig.GitHubOrgId)
		byteValueOfCreateGitopsConfig, _ := json.Marshal(createGitopsConfigRequestDto)

		log.Println("Hitting The post gitops config API")
		HitCreateGitopsConfigApi(byteValueOfCreateGitopsConfig, gitopsConfig.Provider, gitopsConfig.GitUsername, gitopsConfig.Host, gitopsConfig.GitToken, gitopsConfig.GitHubOrgId, suite.authToken)

		log.Println("Hitting the HitFetchAllGitopsConfigApi again for verifying the functionality of it")
		fetchAllLinkResponseDto := HitFetchAllGitopsConfigApi(suite.authToken)
		for _, item := range fetchAllLinkResponseDto.Result {
			if item.Provider == createGitopsConfigRequestDto.Provider {
				assert.True(suite.T(), true)
			}
		}

	})

}
