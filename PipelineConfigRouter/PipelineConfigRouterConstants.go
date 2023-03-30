package PipelineConfigRouter

import "automation-suite/PipelineConfigRouter/ResponseDTOs"

const (
	SaveAppCiPipelineApiUrl                   string = "/orchestrator/app/ci-pipeline"
	SaveAppCiPipelineApi                      string = "SaveAppCiPipelineApi"
	CreateAppApiUrl                           string = "/orchestrator/app"
	CreateAppApi                              string = "CreateAppApi"
	DeleteAppApi                              string = "DeleteAppApi"
	CreateAppMaterialApiUrl                   string = "/orchestrator/app/material"
	CreateAppMaterialApi                      string = "CreateAppMaterialApi"
	DeleteAppMaterialApi                      string = "DeleteAppMaterialApi"
	GetAppDetailsApiUrl                       string = "/orchestrator/app/get"
	GetCiPipelineViaIdApi                     string = "GetCiPipelineViaIdApi"
	GetAppDetailsApi                          string = "FetchAppGetApi"
	GetCiPipelineViaIdApiUrl                  string = "/orchestrator/app/ci-pipeline/"
	GetContainerRegistryApi                   string = "GetContainerRegistryApi"
	PipelineRouterBaseApiUrl                  string = "/orchestrator/app/"
	GetChartReferenceViaAppIdApi              string = "GetChartReferenceViaAppIdApi"
	GetChartReferenceViaAppIdApiUrl           string = "/orchestrator/chartref/autocomplete/"
	GetAppTemplateViaAppIdAndChartRefIdApi    string = "GetAppTemplateViaAppIdAndChartRefIdApi"
	GetAppTemplateViaAppIdAndChartRefIdApiUrl string = "/orchestrator/app/template/"
	GetCdPipelineStrategiesApiUrl             string = "/orchestrator/app/cd-pipeline/strategies/"
	GetCdPipelineStrategiesApi                string = "GetCdPipelineStrategiesApi"
	GetPipelineSuggestedCICDApiUrl            string = "/orchestrator/app/pipeline/suggest/"
	GetPipelineSuggestedCICDApi               string = "GetPipelineSuggestedCICDApi"
	GetEnvAutocompleteApiUrl                  string = "/orchestrator/env/autocomplete"
	GetEnvAutocompleteApi                     string = "GetEnvAutocompleteApi"
	SaveDeploymentTemplateAPiUrl              string = "/orchestrator/app/template"
	SaveDeploymentTemplateApi                 string = "SaveDeploymentTemplateApi"
	PatchCiPipelinesApiUrl                    string = "/orchestrator/app/ci-pipeline/patch"
	PatchCiPipelinesApi                       string = "PatchCiPipelinesApi"
	GetWorkflowDetailsApi                     string = "GetWorkflowDetailsApi"
	GetWorkflowApiUrl                         string = "/orchestrator/app/app-wf/"
	FetchAllAppWorkflowApi                    string = "FetchAllAppWorkflowApi"
	FetchSuggestedCiPipelineNameApi           string = "FetchSuggestedCiPipelineNameApi"
	SaveCdPipelineApiUrl                      string = "/orchestrator/app/cd-pipeline"
	SaveCdPipelineApi                         string = "SaveCdPipelineApi"
	Automatic                                 string = "AUTOMATIC"
	Manual                                    string = "MANUAL"
	DeleteCdPipelineApiUrl                    string = "/orchestrator/app/cd-pipeline/patch"
	ForceDeleteCdPipelineApiUrl               string = "/orchestrator/app/cd-pipeline/patch?force=true"
	DeleteCdPipelineApi                       string = "DeleteCdPipelineApi"
	GetAppCdPipelineApiUrl                    string = "/orchestrator/app/cd-pipeline/"
	GetAppCdPipelineApi                       string = "GetAppCdPipelineApi"
	GetWorkflowStatusApi                      string = "GetWorkflowStatusApi"
	GetWorkflowStatusApiUrl                   string = "/orchestrator/app/workflow/status/"
	GetCiPipelineMaterialApi                  string = "GetCiPipelineMaterialApi"
	GetCiPipelineBaseUrl                      string = "/orchestrator/app/ci-pipeline"
	TriggerCiPipelineApiUrl                   string = GetCiPipelineBaseUrl + "/trigger"
	TriggerCiPipelineApi                      string = "TriggerCiPipelineApi"
	UpdateAppMaterial                         string = "UpdateAppMaterial"
	GetAppListForAutocompleteApi              string = "GetAppListForAutocompleteApi"
	GetAppListForAutocompleteApiUrl           string = "/orchestrator/app/autocomplete"
	GetAppListByTeamIdsApi                    string = "GetAppListByTeamIdsApi"
	GetAppListByTeamIdsApiUrl                 string = "/orchestrator/app/min"
	FindAppsByTeamIdApiUrl                    string = "/orchestrator/app/team/by-id/"
	FindAppsByTeamNameApiUrl                  string = "/orchestrator/app/team/by-name/"
	FetchMaterialsApiUrl                      string = "/orchestrator/app/ci-pipeline/"
	FetchMaterialsApi                         string = "FetchMaterialsApi"
	GetCiPipelineMinApi                       string = "GetCiPipelineMinApi"
	RefreshMaterialsApiUrl                    string = "/orchestrator/app/ci-pipeline/refresh-material/"
	RefreshMaterialsApi                       string = "RefreshMaterialsApi"
	GetAppDeploymentStatusTimelineApi         string = "GetAppDeploymentStatusTimelineApi"
	GetAppDeploymentStatusTimelineApiUrl      string = "/orchestrator/app/deployment-status/timeline/"
	GitListAutocompleteApi                    string = "GitListAutocompleteApi"

	TIMELINE_STATUS_DEPLOYMENT_INITIATED   ResponseDTOs.TimelineStatus = "DEPLOYMENT_INITIATED"
	TIMELINE_STATUS_GIT_COMMIT             ResponseDTOs.TimelineStatus = "GIT_COMMIT"
	TIMELINE_STATUS_GIT_COMMIT_FAILED      ResponseDTOs.TimelineStatus = "GIT_COMMIT_FAILED"
	TIMELINE_STATUS_KUBECTL_APPLY_STARTED  ResponseDTOs.TimelineStatus = "KUBECTL_APPLY_STARTED"
	TIMELINE_STATUS_KUBECTL_APPLY_SYNCED   ResponseDTOs.TimelineStatus = "KUBECTL_APPLY_SYNCED"
	TIMELINE_STATUS_APP_HEALTHY            ResponseDTOs.TimelineStatus = "HEALTHY"
	TIMELINE_STATUS_DEPLOYMENT_FAILED      ResponseDTOs.TimelineStatus = "FAILED"
	TIMELINE_STATUS_FETCH_TIMED_OUT        ResponseDTOs.TimelineStatus = "TIMED_OUT"
	TIMELINE_STATUS_UNABLE_TO_FETCH_STATUS ResponseDTOs.TimelineStatus = "UNABLE_TO_FETCH_STATUS"
	TIMELINE_STATUS_DEPLOYMENT_SUPERSEDED  ResponseDTOs.TimelineStatus = "DEPLOYMENT_SUPERSEDED"
	TIMELINE_RESOURCE_STAGE_KUBECTL_APPLY  string                      = "KUBECTL_APPLY"
)
