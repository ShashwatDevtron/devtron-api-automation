package GitProviderRouter

import (
	"automation-suite/GitProviderRouter/RequestDTOs"
	"automation-suite/GitProviderRouter/ResponseDTOs"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strconv"
)

type StructGitProRouter struct {
	getGitProviderResponse     ResponseDTOs.GetGitProviderResponseDto
	saveGitProviderResponseDto ResponseDTOs.SaveGitProviderResponseDto
	gitRequestDTOs             RequestDTOs.SaveGitProviderRequestDTO
	deleteGitProviderResponse  ResponseDTOs.DeleteGitProviderResponse
	updateGitProviderResponse  ResponseDTOs.SaveGitProviderResponseDto
	GetGitProviderResponseById ResponseDTOs.GetGitProviderResponseById
}

func (structGitRegRouter StructGitProRouter) UnmarshalGivenResponseBody(response []byte, apiName string) StructGitProRouter {
	switch apiName {
	case GetGitProvider:
		json.Unmarshal(response, &structGitRegRouter.getGitProviderResponse)
	case SaveGitProviderApi:
		json.Unmarshal(response, &structGitRegRouter.saveGitProviderResponseDto)
	case UpdateGitProvider:
		json.Unmarshal(response, &structGitRegRouter.saveGitProviderResponseDto)
	case DeleteGitProvider:
		json.Unmarshal(response, &structGitRegRouter.deleteGitProviderResponse)
	case GetGitProviderById:
		json.Unmarshal(response, &structGitRegRouter.GetGitProviderResponseById)
	}

	return structGitRegRouter
}

func HitGetGitProviderApi(authToken string) ResponseDTOs.GetGitProviderResponseDto {
	resp, err := Base.MakeApiCall(GitProviderApiUrl, http.MethodGet, "", nil, authToken)
	Base.HandleError(err, GetGitProvider)

	structGitProRouter := StructGitProRouter{}
	githubProRouter := structGitProRouter.UnmarshalGivenResponseBody(resp.Body(), GetGitProvider)
	return githubProRouter.getGitProviderResponse
}

func HitGetGitProviderByIdApi(appId int, authToken string) ResponseDTOs.GetGitProviderResponseById {
	appIdStr := strconv.Itoa(appId)
	resp, err := Base.MakeApiCall(GitProviderApiUrl+"/"+appIdStr, http.MethodGet, "", nil, authToken)
	Base.HandleError(err, GetGitProviderById)

	structGitProRouter := StructGitProRouter{}
	githubProRouter := structGitProRouter.UnmarshalGivenResponseBody(resp.Body(), GetGitProviderById)
	return githubProRouter.GetGitProviderResponseById
}

func GetGitProviderRequestDto(GitRegHostId int, authMode string) RequestDTOs.SaveGitProviderRequestDTO {
	var saveGitProviderRequestDto RequestDTOs.SaveGitProviderRequestDTO
	envConf := Base.ReadBaseEnvConfig()
	file := Base.ReadAnyJsonFile(envConf.ClassCredentialsFile)
	saveGitProviderRequestDto.Id = 0
	saveGitProviderRequestDto.GitHostId = GitRegHostId
	saveGitProviderRequestDto.Active = true
	saveGitProviderRequestDto.Url = file.GitHubProjectUrl
	saveGitProviderRequestDto.AuthMode = authMode
	if authMode == "SSH" {
		saveGitProviderRequestDto.SshPrivateKey = file.SshPrivateKey
	}
	if authMode == "USERNAME_PASSWORD" {
		saveGitProviderRequestDto.Password = file.GitAccountPassword
	}
	saveGitProviderRequestDto.Name = file.GitProviderName
	return saveGitProviderRequestDto
}

func HitSaveGitProviderApi(payloadOfApi []byte, authToken string) ResponseDTOs.SaveGitProviderResponseDto {
	resp, err := Base.MakeApiCall(GitProviderApiUrl, http.MethodPost, string(payloadOfApi), nil, authToken)
	Base.HandleError(err, SaveGitProviderApi)

	structGitProRouter := StructGitProRouter{}
	githubProRouter := structGitProRouter.UnmarshalGivenResponseBody(resp.Body(), SaveGitProviderApi)
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
	saveGitProviderRequestDto.Name = file.GitProviderUpdatedName
	byteValueOfStruct, _ := json.Marshal(saveGitProviderRequestDto)
	return byteValueOfStruct
}

func HitUpdateGitProviderApi(byteValueOfStruct []byte, authToken string) ResponseDTOs.SaveGitProviderResponseDto {
	resp, err := Base.MakeApiCall(GitProviderApiUrl, http.MethodPut, string(byteValueOfStruct), nil, authToken)
	Base.HandleError(err, UpdateGitProvider)

	structGitProRouter := StructGitProRouter{}
	gitProRouter := structGitProRouter.UnmarshalGivenResponseBody(resp.Body(), UpdateGitProvider)
	return gitProRouter.saveGitProviderResponseDto
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
	resp, err := Base.MakeApiCall(GitProviderApiUrl, http.MethodDelete, string(byteValueOfStruct), nil, authToken)
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
