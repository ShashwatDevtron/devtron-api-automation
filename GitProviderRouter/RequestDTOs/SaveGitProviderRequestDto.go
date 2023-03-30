package RequestDTOs

type SaveGitProviderRequestDTO struct {
	Id            int    `json:"id,omitempty" validate:"number"`
	GitHostId     int    `json:"gitHostId"`
	Active        bool   `json:"active"`
	AuthMode      string `json:"authMode,omitempty" validate:"required"`
	Name          string `json:"name,omitempty" validate:"required"`
	Url           string `json:"url,omitempty"`
	UserName      string `json:"userName,omitempty"`
	Password      string `json:"password,omitempty"`
	SshPrivateKey string `json:"sshPrivateKey,omitempty"`
	AccessToken   string `json:"accessToken,omitempty"`
}
