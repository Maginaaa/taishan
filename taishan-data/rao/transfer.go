package rao

const (
	TypePreSceneExport = iota + 1
)

type TransferInfo struct {
	Type      int    `json:"type"`
	End       bool   `json:"end"`
	MachineIP string `json:"machine_ip"`
	ReportID  int32  `json:"report_id"`
	SceneID   int32  `json:"scene_id"`
	Data      map[string]string
}
