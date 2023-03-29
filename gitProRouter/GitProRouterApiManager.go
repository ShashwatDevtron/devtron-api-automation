package gitProRouter

import (
	"automation-suite/gitProRouter/RequestDTOs"
	"automation-suite/gitProRouter/ResponseDTOs"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"net/http"
)

type StructGitProRouter struct {
	saveGitProviderResponseDto ResponseDTOs.SaveGitProviderResponseDto
	gitRequestDTOs             RequestDTOs.SaveGitProviderRequestDTO
	deleteGitProviderResponse  ResponseDTOs.DeleteGitProviderResponse
	updateGitProviderResponse  ResponseDTOs.SaveGitProviderResponseDto
}

func (structGitRegRouter StructGitProRouter) UnmarshalGivenResponseBody(response []byte, apiName string) StructGitProRouter {
	switch apiName {
	case SaveGitProvideApi:
		json.Unmarshal(response, &structGitRegRouter.saveGitProviderResponseDto)
	}
	return structGitRegRouter
}

func GetGitProviderRequestDto() RequestDTOs.SaveGitProviderRequestDTO {
	var saveGitProviderRequestDto RequestDTOs.SaveGitProviderRequestDTO
	envConf := Base.ReadBaseEnvConfig()
	file := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)
	saveGitProviderRequestDto.Id = 0
	saveGitProviderRequestDto.GitHostId = file.GitHubOrgId
	saveGitProviderRequestDto.Active = true
	saveGitProviderRequestDto.Url = file.GitHubProjectUrl
	saveGitProviderRequestDto.AuthMode = file.AuthMode
	saveGitProviderRequestDto.Name = file.Name
	return saveGitProviderRequestDto
}

func HitSaveGitProviderApi(payloadOfApi []byte, authToken string) ResponseDTOs.SaveGitProviderResponseDto {
	resp, err := Base.MakeApiCall(SaveGitProviderApiUrl, http.MethodPost, string(payloadOfApi), nil, authToken)
	Base.HandleError(err, SaveGitProvideApi)

	structDockerRegRouter := StructGitProRouter{}
	githubProRouter := structDockerRegRouter.UnmarshalGivenResponseBody(resp.Body(), SaveGitProvideApi)
	return githubProRouter.saveGitProviderResponseDto
}
func GetPayLoadForUpdateGitProviderAPI(id int, gitHostId int) []byte {
	var saveGitProviderRequestDto RequestDTOs.SaveGitProviderRequestDTO
	envConf := Base.ReadBaseEnvConfig()
	file := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)
	saveGitProviderRequestDto.Id = id
	saveGitProviderRequestDto.GitHostId = gitHostId
	saveGitProviderRequestDto.Active = true
	saveGitProviderRequestDto.Url = file.GitHubProjectUrl
	saveGitProviderRequestDto.AuthMode = file.AuthMode
	saveGitProviderRequestDto.Name = file.Name
	byteValueOfStruct, _ := json.Marshal(saveGitProviderRequestDto)
	return byteValueOfStruct
}

func HitUpdateGitProviderApi(byteValueOfStruct []byte, authToken string) ResponseDTOs.SaveGitProviderResponseDto {
	resp, err := Base.MakeApiCall(SaveGitProviderApiUrl, http.MethodPut, string(byteValueOfStruct), nil, authToken)
	Base.HandleError(err, UpdateGitProvider)

	structGitProRouter := StructGitProRouter{}
	gitProRouter := structGitProRouter.UnmarshalGivenResponseBody(resp.Body(), UpdateGitProvider)
	return gitProRouter.updateGitProviderResponse
}

func GetPayLoadForDeleteGitProviderAPI(id int, gitHostId int, url string, authMode string, name string) []byte {
	var saveGitProviderRequestDto RequestDTOs.SaveGitProviderRequestDTO
	saveGitProviderRequestDto.Id = id
	saveGitProviderRequestDto.GitHostId = gitHostId
	saveGitProviderRequestDto.Active = true
	saveGitProviderRequestDto.Url = url
	saveGitProviderRequestDto.AuthMode = authMode
	saveGitProviderRequestDto.Name = name
	byteValueOfStruct, _ := json.Marshal(saveGitProviderRequestDto)
	return byteValueOfStruct
}

func HitDeleteGitProviderApi(byteValueOfStruct []byte, authToken string) ResponseDTOs.DeleteGitProviderResponse {
	resp, err := Base.MakeApiCall(SaveGitProviderApiUrl, http.MethodDelete, string(byteValueOfStruct), nil, authToken)
	Base.HandleError(err, DeleteGitProvider)

	structGitProRouter := StructGitProRouter{}
	gitProRouter := structGitProRouter.UnmarshalGivenResponseBody(resp.Body(), DeleteGitProvider)
	return gitProRouter.deleteGitProviderResponse
}

type GitProRouterTestSuite struct {
	suite.Suite
	authToken string
}

func (suite *GitProRouterTestSuite) SetupSuite() {
	suite.authToken = Base.GetAuthToken()
}
