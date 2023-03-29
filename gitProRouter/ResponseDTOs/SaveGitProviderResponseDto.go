package ResponseDTOs

import (
	"automation-suite/gitProRouter/RequestDTOs"
	Base "automation-suite/testUtils"
)

type SaveGitProviderResponseDto struct {
	Code   int                                   `json:"code"`
	Status string                                `json:"status"`
	Result RequestDTOs.SaveGitProviderRequestDTO `json:"result"`
	Errors []Base.Errors                         `json:"errors"`
}

type DeleteGitProviderResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Result string `json:"result"`
}
