package GitopsConfigRouter

import (
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
)

func (suite *GitOpsRouterTestSuite) TestClassA5UpdateGitopsConfig() {
	suite.Run("A=7=UpdateGitopsConfigWithValidPayload", func() {
		envConf := Base.ReadBaseEnvConfig()
		gitopsConfig := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)

		createGitopsConfigRequestDto := GetGitopsConfigRequestDto(gitopsConfig.Provider, gitopsConfig.GitUsername, gitopsConfig.Host, gitopsConfig.GitToken, gitopsConfig.GitHubOrgId)
		byteValueOfCreateGitopsConfig, _ := json.Marshal(createGitopsConfigRequestDto)
		ValidatedResponse := HitValidateGitopsConfigApi(byteValueOfCreateGitopsConfig, gitopsConfig.Provider, gitopsConfig.GitUsername, gitopsConfig.Host, gitopsConfig.GitToken, gitopsConfig.GitHubOrgId, suite.authToken)
		log.Println("Validating the Response of the Validate Gitops Config API...")
		assert.Equal(suite.T(), 200, ValidatedResponse.Code)
		assert.Equal(suite.T(), "Create Repo", ValidatedResponse.Result.SuccessfulStages[0])
	})

}
