package PipelineConfigRouter

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
	GetContainerRegistryApiUrl                string = "/orchestrator/app/"
	GetChartReferenceViaAppIdApi              string = "GetChartReferenceViaAppIdApi"
	GetChartReferenceViaAppIdApiUrl           string = "/orchestrator/chartref/autocomplete/"
	GetAppTemplateViaAppIdAndChartRefIdApi    string = "GetAppTemplateViaAppIdAndChartRefIdApi"
	GetAppTemplateViaAppIdAndChartRefIdApiUrl string = "/orchestrator/app/template/"
	GetCdPipelineStrategiesApiUrl             string = "/orchestrator/app/cd-pipeline/strategies/"
	GetCdPipelineStrategiesApi                string = "GetCdPipelineStrategiesApi"
	GetPipelineSuggestedCDApiUrl              string = "/orchestrator/app/pipeline/suggest/cd/"
	GetPipelineSuggestedCDApi                 string = "GetPipelineSuggestedCDApi"
	GetAllEnvironmentDetailsApiUrl            string = "/orchestrator/env/autocomplete?auth=true"
	GetAllEnvironmentDetailsApi               string = "GetAllEnvironmentDetailsApiUrl"
	SaveDeploymentTemplateAPiUrl              string = "/orchestrator/app/template"
	SaveDeploymentTemplateApi                 string = "SaveDeploymentTemplateApi"
	CreateWorkflowApiUrl                      string = "/orchestrator/app/ci-pipeline/patch"
	CreateWorkflowApi                         string = "CreateWorkflowApi"
	GetWorkflowDetailsApi                     string = "GetWorkflowDetailsApi"
	DeleteWorkflowApiUrl                      string = "/orchestrator/app/app-wf/"
	FetchSuggestedCiPipelineNameApiUrl        string = "/orchestrator/app/pipeline/suggest/ci/"
	FetchSuggestedCiPipelineNameApi           string = "FetchSuggestedCiPipelineNameApi"
	SaveCdPipelineApiUrl                      string = "/orchestrator/app/cd-pipeline"
	SaveCdPipelineApi                         string = "SaveCdPipelineApi"
	Automatic                                 string = "AUTOMATIC"
	Manual                                    string = "MANUAL"
	DeleteCdPipelineApiUrl                    string = "/orchestrator/app/cd-pipeline/patch"
	DeleteCdPipelineApi                       string = "DeleteCdPipelineApi"
)
