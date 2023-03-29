package gitProRouter

import (
	"automation-suite/gitProRouter/ResponseDTOs"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
)

func (suite *GitProRouterTestSuite) TestClassA6GetApp() {

	var byteValueOfSaveGitProvider []byte
	var saveGitProviderResponseDto ResponseDTOs.SaveGitProviderResponseDto

	suite.Run("A=1=SaveGitProviderWithValidPayload", func() {
		saveGitProviderRequestDto := GetGitProviderRequestDto()
		byteValueOfSaveGitProvider, _ = json.Marshal(saveGitProviderRequestDto)

		log.Println("Hitting The post git provider API")
		saveGitProviderResponseDto = HitSaveGitProviderApi(byteValueOfSaveGitProvider, suite.authToken)

		log.Println("Validating the Response of the save git provider API...")
		assert.Equal(suite.T(), saveGitProviderRequestDto.Name, saveGitProviderResponseDto.Result.Name)
	})

	suite.Run("A=2=SaveGitProviderWithExistingName", func() {
		saveGitProviderRequestPayload := GetGitProviderRequestDto()
		byteValueOfGitProviderPayload, _ := json.Marshal(saveGitProviderRequestPayload)

		log.Println("Hitting The save git provider Api First time")
		saveGitProviderResponse := HitSaveGitProviderApi(byteValueOfGitProviderPayload, suite.authToken)

		log.Println("Hitting The save git provider Api second time with existing registry name")
		finalApiResponse := HitSaveGitProviderApi(byteValueOfGitProviderPayload, suite.authToken)

		log.Println("Validating the Response of the save git provider  API...")
		assert.Equal(suite.T(), finalApiResponse.Status, "Internal Server Error")
		assert.Equal(suite.T(), finalApiResponse.Status, "Internal Server Error")

		log.Println("getting payload for git provider API")
		byteValueOfDeleteGitProvider := GetPayLoadForDeleteGitProviderAPI(saveGitProviderResponse.Result.Id, saveGitProviderResponse.Result.GitHostId, saveGitProviderResponse.Result.Url, saveGitProviderResponse.Result.AuthMode, saveGitProviderResponse.Result.Name)
		log.Println("Hitting the Delete team API for Removing the data created via automation")
		HitDeleteGitProviderApi(byteValueOfDeleteGitProvider, suite.authToken)
	})

	log.Println("getting payload for Delete git provider API")
	byteValueOfDeleteDockerRegistry := GetPayLoadForDeleteGitProviderAPI(saveGitProviderResponseDto.Result.Id, saveGitProviderResponseDto.Result.GitHostId, saveGitProviderResponseDto.Result.Url, saveGitProviderResponseDto.Result.AuthMode, saveGitProviderResponseDto.Result.Name)
	log.Println("Hitting the Delete team API for Removing the data created via automation")
	HitDeleteGitProviderApi(byteValueOfDeleteDockerRegistry, suite.authToken)
}
