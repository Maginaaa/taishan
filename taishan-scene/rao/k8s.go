package rao

type DeploymentQueryReq struct {
	Namespace string `json:"namespace"`
}

type DeploymentInfo struct {
	Namespace      string `json:"namespace"`
	DeploymentName string `json:"deployment_name"`
	Replicas       int32  `json:"replicas"`
}
