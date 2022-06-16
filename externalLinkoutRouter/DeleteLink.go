package externalLinkoutRouter

import (
	"automation-suite/testUtils"
	"log"

	"github.com/stretchr/testify/assert"
)

func (suite *ExternalLinkOutRouterTestSuite) TestClassA2DeleteExternalLink() {
	/*suite.Run("A=1=DeleteExternalLinkoutWithValidId", func() {
		log.Println("Fetching links before creating new")
		getAllExternalLinksResponseDto := HitFetchAllLinkApi(suite.authToken)
		noOfLinks := len(getAllExternalLinksResponseDto.Result)
		createLinkRequestDto := GetSaveLinkRequestDto(1, nil)
		byteValueOfCreateLink, _ := json.Marshal(createLinkRequestDto)
		log.Println("Hitting The Save Link API")
		HitCreateLinkApi(byteValueOfCreateLink, suite.authToken)
		getAllExternalLinksAgainResponseDto := HitFetchAllLinkApi(suite.authToken)
		noOfLinksAfterCreation := len(getAllExternalLinksAgainResponseDto.Result)
		log.Println("Hitting the Delete link API for Removing the data created via automation")
		externalLinkdeleteResponseDto := HitDeleteLinkApi(getAllExternalLinksAgainResponseDto.Result[noOfLinksAfterCreation-1].Id, suite.authToken)
		assert.Equal(suite.T(), true, externalLinkdeleteResponseDto.Result.Success)
		getAllExternalLinksAgainToCheckResponseDto := HitFetchAllLinkApi(suite.authToken)
		noOfLinksToCheck := len(getAllExternalLinksAgainToCheckResponseDto.Result)
		assert.Equal(suite.T(), noOfLinks, noOfLinksToCheck)
	})*/
	suite.Run("A=2=DeleteExternalLinkoutWithInValidId", func() {
		log.Println("Hitting the Delete link API for Removing the data created via automation")
		externalLinkdeleteResponseDto := HitDeleteLinkApi(testUtils.GetRandomNumberOf9Digit(), suite.authToken)
		assert.Equal(suite.T(), 404, externalLinkdeleteResponseDto.Code)
		assert.Equal(suite.T(), "pg: no rows in result set", externalLinkdeleteResponseDto.Errors[0].UserMessage)

	})
}