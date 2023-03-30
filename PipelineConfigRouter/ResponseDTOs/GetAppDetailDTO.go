package ResponseDTOs

import "time"

type GetAppDetail struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Result struct {
		AppId            int    `json:"appId"`
		AppName          string `json:"appName"`
		EnvironmentId    int    `json:"environmentId"`
		EnvironmentName  string `json:"environmentName"`
		Namespace        string `json:"namespace"`
		LastDeployedTime string `json:"lastDeployedTime"`
		LastDeployedBy   string `json:"lastDeployedBy"`
		MaterialInfo     []struct {
			Author       string    `json:"author"`
			Branch       string    `json:"branch"`
			Message      string    `json:"message"`
			ModifiedTime time.Time `json:"modifiedTime"`
			Revision     string    `json:"revision"`
			Url          string    `json:"url"`
			WebhookData  string    `json:"webhookData"`
		} `json:"materialInfo"`
		ReleaseVersion             string      `json:"releaseVersion"`
		DataSource                 string      `json:"dataSource"`
		LastDeployedPipeline       string      `json:"lastDeployedPipeline"`
		Deprecated                 bool        `json:"deprecated"`
		K8SVersion                 string      `json:"k8sVersion"`
		CiArtifactId               int         `json:"ciArtifactId"`
		ClusterId                  int         `json:"clusterId"`
		DeploymentAppType          string      `json:"deploymentAppType"`
		ExternalCi                 bool        `json:"externalCi"`
		ClusterName                string      `json:"clusterName"`
		DockerRegistryId           string      `json:"dockerRegistryId"`
		IpsAccessProvided          bool        `json:"ipsAccessProvided"`
		DeploymentAppDeleteRequest bool        `json:"deploymentAppDeleteRequest"`
		InstanceDetail             interface{} `json:"instanceDetail"`
		OtherEnvironment           []struct {
			AppStatus                  string `json:"appStatus"`
			EnvironmentId              int    `json:"environmentId"`
			EnvironmentName            string `json:"environmentName"`
			AppMetrics                 bool   `json:"appMetrics"`
			InfraMetrics               bool   `json:"infraMetrics"`
			Prod                       bool   `json:"prod"`
			ChartRefId                 int    `json:"chartRefId"`
			LastDeployed               string `json:"lastDeployed"`
			DeploymentAppDeleteRequest bool   `json:"deploymentAppDeleteRequest"`
		} `json:"otherEnvironment"`
		ResourceTree struct {
			Conditions interface{} `json:"conditions"`
			Hosts      []struct {
				Name          string `json:"name"`
				ResourcesInfo []struct {
					Capacity             int64  `json:"capacity"`
					RequestedByApp       int64  `json:"requestedByApp"`
					RequestedByNeighbors int64  `json:"requestedByNeighbors"`
					ResourceName         string `json:"resourceName"`
				} `json:"resourcesInfo"`
				SystemInfo struct {
					Architecture            string `json:"architecture"`
					BootID                  string `json:"bootID"`
					ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
					KernelVersion           string `json:"kernelVersion"`
					KubeProxyVersion        string `json:"kubeProxyVersion"`
					KubeletVersion          string `json:"kubeletVersion"`
					MachineID               string `json:"machineID"`
					OperatingSystem         string `json:"operatingSystem"`
					OsImage                 string `json:"osImage"`
					SystemUUID              string `json:"systemUUID"`
				} `json:"systemInfo"`
			} `json:"hosts"`
			NewGenerationReplicaSets []string `json:"newGenerationReplicaSets"`
			Nodes                    []struct {
				CreatedAt       time.Time `json:"createdAt"`
				Kind            string    `json:"kind"`
				Name            string    `json:"name"`
				Namespace       string    `json:"namespace"`
				ResourceVersion string    `json:"resourceVersion"`
				Uid             string    `json:"uid"`
				Version         string    `json:"version"`
				ParentRefs      []struct {
					Kind      string `json:"kind"`
					Name      string `json:"name"`
					Namespace string `json:"namespace"`
					Uid       string `json:"uid"`
					Group     string `json:"group,omitempty"`
				} `json:"parentRefs,omitempty"`
				Health struct {
					Status  string `json:"status"`
					Message string `json:"message,omitempty"`
				} `json:"health,omitempty"`
				Images []string `json:"images,omitempty"`
				Info   []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"info,omitempty"`
				NetworkingInfo struct {
					Labels struct {
						App                     string `json:"app"`
						AppId                   string `json:"appId"`
						EnvId                   string `json:"envId"`
						Release                 string `json:"release"`
						RolloutsPodTemplateHash string `json:"rollouts-pod-template-hash"`
					} `json:"labels,omitempty"`
					TargetLabels struct {
						App string `json:"app"`
					} `json:"targetLabels,omitempty"`
				} `json:"networkingInfo,omitempty"`
				Group string `json:"group,omitempty"`
			} `json:"nodes"`
			PodMetadata []struct {
				Containers     []string    `json:"containers"`
				InitContainers interface{} `json:"initContainers"`
				IsNew          bool        `json:"isNew"`
				Name           string      `json:"name"`
				Uid            string      `json:"uid"`
			} `json:"podMetadata"`
			RevisionHash string `json:"revisionHash"`
			Status       string `json:"status"`
		} `json:"resourceTree"`
	} `json:"result"`
}
