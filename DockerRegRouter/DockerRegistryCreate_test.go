package DockerRegRouter

import (
	"automation-suite/regressionTestSuite"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
)

func (suite *DockerRegRouter) TestSaveDockerRegistryWithValidPayload() {
	saveDockerRegistryRequestDto := GetDockerRegistryRequestDto(false, "", "", "", "", false, "", "")
	byteValueOfSaveDockerRegistry, _ := json.Marshal(saveDockerRegistryRequestDto)

	log.Println("Hitting The post Docker registry API")
	saveDockerRegistryResponseDto := HitSaveDockerRegistryApi(false, byteValueOfSaveDockerRegistry, "", "", "", "", "", "", false, suite.authToken)

	log.Println("Validating the Response of the save docker registry API...")
	assert.Equal(suite.T(), saveDockerRegistryRequestDto.Id, saveDockerRegistryResponseDto.Result.Id)

	log.Println("getting payload for Delete Team API")
	byteValueOfDeleteDockerRegistry := GetPayLoadForDeleteDockerRegistryAPI(saveDockerRegistryResponseDto.Result.Id, saveDockerRegistryResponseDto.Result.PluginId, saveDockerRegistryResponseDto.Result.RegistryUrl, saveDockerRegistryResponseDto.Result.RegistryType, saveDockerRegistryResponseDto.Result.Username, saveDockerRegistryResponseDto.Result.Password, saveDockerRegistryResponseDto.Result.IsDefault)
	log.Println("Hitting the Delete team API for Removing the data created via automation")
	HitDeleteDockerRegistryApi(byteValueOfDeleteDockerRegistry, suite.authToken)
}
func (suite *DockerRegRouter) TestSaveDockerRegistryWithPreviousId() {
	saveDockerRegistryRequestDto := regressionTestSuite.GetDockerRegistryRequestDto(false, "", "", "", "", false, "", "")
	byteValueOfSaveDockerRegistry, _ := json.Marshal(saveDockerRegistryRequestDto)

	log.Println("Hitting The save docker registry API")
	saveDockerRegistryResponseDto := HitSaveDockerRegistryApi(false, byteValueOfSaveDockerRegistry, "", "", "", "", "", "", false, suite.authToken)

	log.Println("Hitting HitSaveDockerRegistryApi with same payload again")
	saveDockerRegistryOnceAgainRequestDto := GetDockerRegistryRequestDto(true, saveDockerRegistryResponseDto.Result.Id, saveDockerRegistryResponseDto.Result.PluginId, saveDockerRegistryResponseDto.Result.RegistryType, saveDockerRegistryResponseDto.Result.RegistryUrl, saveDockerRegistryResponseDto.Result.IsDefault, saveDockerRegistryResponseDto.Result.Username, saveDockerRegistryResponseDto.Result.Password)
	byteValueOfSaveDockerRegistryOnceAgain, _ := json.Marshal(saveDockerRegistryOnceAgainRequestDto)
	saveDockerRegistryOnceAgainResponseDto := HitSaveDockerRegistryApi(true, byteValueOfSaveDockerRegistryOnceAgain, saveDockerRegistryResponseDto.Result.Id, saveDockerRegistryResponseDto.Result.PluginId, saveDockerRegistryResponseDto.Result.RegistryUrl, saveDockerRegistryResponseDto.Result.RegistryType, saveDockerRegistryResponseDto.Result.Username, saveDockerRegistryResponseDto.Result.Password, saveDockerRegistryResponseDto.Result.IsDefault, suite.authToken)

	log.Println("Validating the Response of the save docker registry API...")
	assert.Equal(suite.T(), "docker registry failed to create in db", saveDockerRegistryOnceAgainResponseDto.Errors[0].InternalMessage)

	log.Println("getting payload for delete docker registry API")
	byteValueOfDeleteDockerRegistry := regressionTestSuite.GetPayLoadForDeleteDockerRegistryAPI(saveDockerRegistryResponseDto.Result.Id, saveDockerRegistryResponseDto.Result.PluginId, saveDockerRegistryResponseDto.Result.RegistryUrl, saveDockerRegistryResponseDto.Result.RegistryType, saveDockerRegistryResponseDto.Result.Username, saveDockerRegistryResponseDto.Result.Password, saveDockerRegistryResponseDto.Result.IsDefault)
	log.Println("Hitting the Delete docker registry API for Removing the data created via automation")
	HitDeleteDockerRegistryApi(byteValueOfDeleteDockerRegistry, suite.authToken)
}
