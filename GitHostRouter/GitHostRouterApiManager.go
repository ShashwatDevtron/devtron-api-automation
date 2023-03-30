package GitHostRouter

import (
	"automation-suite/GitHostRouter/ResponseDTOs"
	Base "automation-suite/testUtils"
	"encoding/json"
	"net/http"
)

type StructGitHostRouter struct {
	getGitHostResponse ResponseDTOs.GetGitHostResponseDto
}

func (structGitRegRouter StructGitHostRouter) UnmarshalGivenResponseBody(response []byte, apiName string) StructGitHostRouter {
	switch apiName {
	case GetGitHost:
		json.Unmarshal(response, &structGitRegRouter.getGitHostResponse)

	}
	return structGitRegRouter
}

func HitGetGitHostApi(authToken string) ResponseDTOs.GetGitHostResponseDto {
	resp, err := Base.MakeApiCall(GitHostApiUrl, http.MethodGet, "", nil, authToken)
	Base.HandleError(err, GetGitHost)

	structGitHostRouter := StructGitHostRouter{}
	githubProRouter := structGitHostRouter.UnmarshalGivenResponseBody(resp.Body(), GetGitHost)
	return githubProRouter.getGitHostResponse
}
