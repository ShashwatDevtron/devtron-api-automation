package UserRouter

const (
	SuperAdmin                 string = "SuperAdmin"
	GroupsAndRoleFilter        string = "GroupsAndRoleFilter"
	GroupsAndRoleFilterDynamic string = "GroupsAndRoleFilterDynamic"
	RoleFilterOnly             string = "RoleFilterOnly"
	GroupsOnly                 string = "GroupsOnly"
	WithHelmAppsOnly           string = "WithHelmAppsOnly"
	WithDevtronAppsOnly        string = "WithDevtronAppsOnly"
	WithChartGroupsOnly        string = "WithChartGroupsOnly"
	WithAllFilter              string = "WithAllFilter"
	CreatUserApiUrl            string = "/orchestrator/user"
	GetUserByIdApiUrl          string = "/orchestrator/user/"
	CreateUserApiUrl           string = "/orchestrator/user"
	UpdateUserApiUrl           string = "/orchestrator/user"
	DeleteUserApiUrl           string = "/orchestrator/user/"
	CreateRoleGroupApiUrl      string = "/orchestrator/user/role/group"
	GetRoleGroupByIdApiUrl     string = "/orchestrator/user/role/group/"
	DeleteRoleGroupByIdApiUrl  string = "/orchestrator/user/role/group/"
	WithDevtronAppsOnlyDynamic string = "WithDevtronAppsOnlyDynamic"
	GLOBALCONFIGURATIONS       string = "globalconfigurations"
	CREATEAPP                  string = "createApp"
	CREATEUSER                 string = "createUser"
	APPLISTFETCH               string = "appListFetch"
	PIPELINECREATE             string = "pipeLineCreate"
	PIPELINEDELETE             string = "pipeLineDelete"
	CREATEAPPMATERIAL          string = "createAppMaterial"
	SAVECDPIPELINE             string = "saveCdPipeline"
	PIPELINEFETCH              string = "pipeLineFetch"
	TRIGGERPIPELINE            string = "triggerPipeline"
	APPDETAILS                 string = "appDetails"
	ENVIRONMENTOVERRIDES       string = "environmentOverrides"
	DELETE                     string = "delete"
	PROJECT                    string = "test29"
	ENV                        string = "envtest29"
	APP                        string = "apptest29"
	ACTION                     string = ""
	ACCESS_TYPE                string = "devtron-app"
	ENTITY                     string = "apps"
)
