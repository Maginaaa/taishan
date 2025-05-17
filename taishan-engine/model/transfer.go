package model

import "encoding/json"

const (
	TypePreSceneExport = iota + 1
)

type TransferInfo[T any] struct {
	Type      int    `json:"type"`
	End       bool   `json:"end"`
	MachineIP string `json:"machine_ip"`
	ReportID  int32  `json:"report_id"`
	SceneID   int32  `json:"scene_id"`
	Data      T
}

func (t *TransferInfo[T]) ToByte() (msg []byte) {
	msg, _ = json.Marshal(t)
	return
}

type ExportDataMap map[string]string
