package model

import (
	"github.com/shirou/gopsutil/load"
)

type HeartBeat struct {
	Name     string        `json:"name"`
	IP       string        `json:"ip"`
	CpuUsage float64       `json:"cpu_usage"`
	CpuLoad  *load.AvgStat `json:"cpu_load"`
	MemInfo  MemInfo       `json:"mem_info"`
	Networks []Network     `json:"networks"`
	//DiskInfos []DiskInfo    `json:"disk_infos"`
	//ServerType        int           `json:"server_type"`
	CreateTime int64 `json:"create_time"`
}

type MemInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

type DiskInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type Network struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
}

type MachineInfoReq struct {
	IPList []string `json:"ip_list"`
}
