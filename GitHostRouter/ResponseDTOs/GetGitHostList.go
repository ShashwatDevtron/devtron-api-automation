package ResponseDTOs

import Base "automation-suite/testUtils"

type GitHostResponse struct {
	Id              int    `json:"id,omitempty" validate:"number"`
	Name            string `json:"name,omitempty" validate:"required"`
	Active          bool   `json:"active"`
	WebhookUrl      string `json:"webhookUrl"`
	WebhookSecret   string `json:"webhookSecret"`
	EventTypeHeader string `json:"eventTypeHeader"`
	SecretHeader    string `json:"secretHeader"`
	SecretValidator string `json:"secretValidator"`
	UserId          int32  `json:"-"`
}

type GetGitHostResponseDto struct {
	Code   int               `json:"code"`
	Status string            `json:"status"`
	Result []GitHostResponse `json:"result"`
	Errors []Base.Errors     `json:"errors"`
}
