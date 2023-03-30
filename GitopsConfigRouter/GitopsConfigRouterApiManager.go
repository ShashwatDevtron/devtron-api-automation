package GitopsConfigRouter

import (
	"automation-suite/GitopsConfigRouter/RequestDTOs"
	"automation-suite/GitopsConfigRouter/ResponseDTOs"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"strconv"
)

type StructGitopsConfigRouter struct {
	createGitopsConfigResponseDto    ResponseDTOs.CreateGitopsConfigResponseDto
	fetchAllGitopsConfigResponseDto  ResponseDTOs.FetchAllGitopsConfigResponseDto
	checkGitopsExistsResponse        ResponseDTOs.CheckGitopsExistsResponse
	updateGitopsConfigResponseDto    ResponseDTOs.UpdateGitopsConfigResponseDto
	fetchGitopsConfigResponseByIdDto ResponseDTOs.FetchGitopsConfigResponseByIdDto
}

func HitGitopsConfigured(authToken string) ResponseDTOs.CheckGitopsExistsResponse {
	resp, err := Base.MakeApiCall(CheckGitopsConfigExistsApiUrl, http.MethodGet, "", nil, authToken)
	Base.HandleError(err, CheckGitopsConfigExistsApi)

	structGitopsConfigRouter := StructGitopsConfigRouter{}
	gitopsConfigRouter := structGitopsConfigRouter.UnmarshalGivenResponseBody(resp.Body(), CheckGitopsConfigExistsApi)
	return gitopsConfigRouter.checkGitopsExistsResponse
}
func HitValidateGitopsConfigApi(payload []byte, provider string, username string, host string, token string, githuborgid string, authToken string) ResponseDTOs.CreateGitopsConfigResponseDto {
	var payloadOfApi string
	if payload != nil {
		payloadOfApi = string(payload)
	} else {
		var createGitopsConfigRequestDto RequestDTOs.CreateGitopsConfigRequestDto
		createGitopsConfigRequestDto.Provider = provider
		createGitopsConfigRequestDto.Username = username
		createGitopsConfigRequestDto.Host = host
		createGitopsConfigRequestDto.Token = token
		createGitopsConfigRequestDto.GitHubOrgId = githuborgid
		createGitopsConfigRequestDto.Active = true
		byteValueOfStruct, _ := json.Marshal(createGitopsConfigRequestDto)
		payloadOfApi = string(byteValueOfStruct)
	}

	resp, err := Base.MakeApiCall(ValidateGitopsConfigApi, http.MethodPost, payloadOfApi, nil, authToken)
	Base.HandleError(err, CreateGitopsConfigApi)

	structGitopsConfigRouter := StructGitopsConfigRouter{}
	gitopsConfigRouter := structGitopsConfigRouter.UnmarshalGivenResponseBody(resp.Body(), CreateGitopsConfigApi)
	return gitopsConfigRouter.createGitopsConfigResponseDto
}

func (structGitopsConfigRouter StructGitopsConfigRouter) UnmarshalGivenResponseBody(response []byte, apiName string) StructGitopsConfigRouter {
	switch apiName {
	case FetchAllGitopsConfigApi:
		json.Unmarshal(response, &structGitopsConfigRouter.fetchAllGitopsConfigResponseDto)
	case CreateGitopsConfigApi:
		json.Unmarshal(response, &structGitopsConfigRouter.createGitopsConfigResponseDto)
	case CheckGitopsConfigExistsApi:
		json.Unmarshal(response, &structGitopsConfigRouter.checkGitopsExistsResponse)
	case UpdateGitopsConfigApi:
		json.Unmarshal(response, &structGitopsConfigRouter.updateGitopsConfigResponseDto)
	case FetchGitopsConfigByIdApi:
		json.Unmarshal(response, &structGitopsConfigRouter.fetchGitopsConfigResponseByIdDto)

	}
	return structGitopsConfigRouter
}

func HitFetchAllGitopsConfigApi(authToken string) ResponseDTOs.FetchAllGitopsConfigResponseDto {
	resp, err := Base.MakeApiCall(SaveGitopsConfigApiUrl, http.MethodGet, "", nil, authToken)
	Base.HandleError(err, FetchAllGitopsConfigApi)

	structGitopsConfigRouter := StructGitopsConfigRouter{}
	gitopsConfigRouter := structGitopsConfigRouter.UnmarshalGivenResponseBody(resp.Body(), FetchAllGitopsConfigApi)
	return gitopsConfigRouter.fetchAllGitopsConfigResponseDto
}
func HitFetchGitopsConfigByIdApi(payload []byte, id int, authToken string) ResponseDTOs.FetchGitopsConfigResponseByIdDto {
	resp, err := Base.MakeApiCall(SaveGitopsConfigApiUrl+"/"+strconv.Itoa(id), http.MethodGet, "", nil, authToken)
	Base.HandleError(err, FetchGitopsConfigByIdApi)

	structGitopsConfigRouter := StructGitopsConfigRouter{}
	gitopsConfigRouter := structGitopsConfigRouter.UnmarshalGivenResponseBody(resp.Body(), FetchAllGitopsConfigApi)
	return gitopsConfigRouter.fetchGitopsConfigResponseByIdDto
}

func GetGitopsConfigRequestDto(provider string, username string, host string, token string, githuborgid string) RequestDTOs.CreateGitopsConfigRequestDto {
	var createGitopsConfigRequestDto RequestDTOs.CreateGitopsConfigRequestDto
	createGitopsConfigRequestDto.Provider = provider
	createGitopsConfigRequestDto.Username = username
	createGitopsConfigRequestDto.Host = host
	createGitopsConfigRequestDto.Token = token
	createGitopsConfigRequestDto.GitHubOrgId = githuborgid
	createGitopsConfigRequestDto.Active = true
	return createGitopsConfigRequestDto
}
func HitCreateGitopsConfigApi(payload []byte, provider string, username string, host string, token string, githuborgid string, authToken string) ResponseDTOs.CreateGitopsConfigResponseDto {
	var payloadOfApi string
	if payload != nil {
		payloadOfApi = string(payload)
	} else {
		var createGitopsConfigRequestDto RequestDTOs.CreateGitopsConfigRequestDto
		createGitopsConfigRequestDto.Provider = provider
		createGitopsConfigRequestDto.Username = username
		createGitopsConfigRequestDto.Host = host
		createGitopsConfigRequestDto.Token = token
		createGitopsConfigRequestDto.GitHubOrgId = githuborgid
		createGitopsConfigRequestDto.Active = true
		byteValueOfStruct, _ := json.Marshal(createGitopsConfigRequestDto)
		payloadOfApi = string(byteValueOfStruct)
	}

	resp, err := Base.MakeApiCall(SaveGitopsConfigApiUrl, http.MethodPost, payloadOfApi, nil, authToken)
	Base.HandleError(err, CreateGitopsConfigApi)

	structGitopsConfigRouter := StructGitopsConfigRouter{}
	gitopsConfigRouter := structGitopsConfigRouter.UnmarshalGivenResponseBody(resp.Body(), CreateGitopsConfigApi)
	return gitopsConfigRouter.createGitopsConfigResponseDto
}

func UpdateGitops(authToken string) RequestDTOs.CreateGitopsConfigRequestDto {
	var createGitopsConfigRequestDto RequestDTOs.CreateGitopsConfigRequestDto
	fetchAllLinkResponseDto := HitFetchAllGitopsConfigApi(authToken)

	log.Println("Checking which is true")
	for _, createGitopsConfigRequestDto = range fetchAllLinkResponseDto.Result {
		if createGitopsConfigRequestDto.Active {
			createGitopsConfigRequestDto.Active = false
			byteValueOfCreateGitopsConfig, _ := json.Marshal(createGitopsConfigRequestDto)
			log.Println("Updating gitops to false")
			HitUpdateGitopsConfigApi(byteValueOfCreateGitopsConfig, authToken)
			createGitopsConfigRequestDto.Active = true
			return createGitopsConfigRequestDto
		}
	}
	return createGitopsConfigRequestDto
}

func HitUpdateGitopsConfigApi(payload []byte, authToken string) ResponseDTOs.UpdateGitopsConfigResponseDto {
	resp, err := Base.MakeApiCall(SaveGitopsConfigApiUrl, http.MethodPut, string(payload), nil, authToken)
	Base.HandleError(err, UpdateGitopsConfigApi)

	structGitopsConfigRouter := StructGitopsConfigRouter{}
	gitopsConfigRouter := structGitopsConfigRouter.UnmarshalGivenResponseBody(resp.Body(), UpdateGitopsConfigApi)
	return gitopsConfigRouter.updateGitopsConfigResponseDto
}

type GitOpsRouterTestSuite struct {
	suite.Suite
	authToken string
}

func (suite *GitOpsRouterTestSuite) SetupSuite() {
	suite.authToken = Base.GetAuthToken()
}
