package rao

type PageResponse struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}

const (
	StopPlan     = 1
	DebugStatus  = 2
	ReportChange = 3
	SceneRelease = 4
)
