package RequestDTOs

import "automation-suite/PipelineConfigRouter"


type SearchDevtronDeletion struct {
	ProjectPayload []byte
	EnvPayLoad     []byte
	DevtronPayload PipelineConfigRouter.CreateAppResponseDto
	RoleGroupId    int
}