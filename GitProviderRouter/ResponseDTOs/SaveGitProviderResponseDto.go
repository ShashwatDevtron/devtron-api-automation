package ResponseDTOs

import (
	"automation-suite/GitProviderRouter/RequestDTOs"
	Base "automation-suite/testUtils"
)

type GetGitProviderResponse struct {
	Id            int    `json:"id,omitempty" validate:"number"`
	Name          string `json:"name,omitempty" validate:"required"`
	Url           string `json:"url,omitempty"`
	UserName      string `json:"userName,omitempty"`
	Password      string `json:"password,omitempty"`
	SshPrivateKey string `json:"sshPrivateKey,omitempty"`
	AccessToken   string `json:"accessToken,omitempty"`
	AuthMode      string `json:"authMode,omitempty" validate:"required"`
	Active        bool   `json:"active"`
	UserId        int32  `json:"-"`
	GitHostId     int    `json:"gitHostId"`
}

type GetGitProviderResponseById struct {
	Code   int                    `json:"code"`
	Status string                 `json:"status"`
	Result GetGitProviderResponse `json:"result"`
	Errors []Base.Errors          `json:"errors"`
}

type GetGitProviderResponseDto struct {
	Code   int                      `json:"code"`
	Status string                   `json:"status"`
	Result []GetGitProviderResponse `json:"result"`
	Errors []Base.Errors            `json:"errors"`
}

type SaveGitProviderResponseDto struct {
	Code   int                                   `json:"code"`
	Status string                                `json:"status"`
	Result RequestDTOs.SaveGitProviderRequestDTO `json:"result"`
	Errors []Base.Errors                         `json:"errors"`
}

type DeleteGitProviderResponse struct {
	Code   int           `json:"code"`
	Status string        `json:"status"`
	Result string        `json:"result"`
	Errors []Base.Errors `json:"errors"`
}
