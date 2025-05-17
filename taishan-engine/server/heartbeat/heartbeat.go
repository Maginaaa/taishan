package heartbeat

import (
	"encoding/json"
	"engine/config"
	"engine/internal/biz/log"
	"engine/middleware"
	"engine/model"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/net"
	gonet "net"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var (
	heartbeat = new(HeartBeat)
	Key       = "TaishanMachineList"
)

func CheckHeartBeat() *HeartBeat {
	heartbeat.Name = GetHostName()
	heartbeat.IP = GetIP()
	heartbeat.CpuUsage = GetCpuUsed()
	heartbeat.MemInfo = GetMemInfo()
	heartbeat.CpuLoad = GetCPULoad()
	heartbeat.Networks = GetNetwork()
	heartbeat.CreateTime = time.Now().Unix()
	return heartbeat
}

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

// GetCpuUsed CPU信息
func GetCpuUsed() float64 {
	percent, _ := cpu.Percent(time.Second, false) // false表示CPU总使用率，true为单核
	return percent[0]
}

// GetCPULoad 负载信息
func GetCPULoad() (info *load.AvgStat) {
	info, _ = load.Avg()
	return
}

// GetMemInfo 内存信息
func GetMemInfo() (memInfoList MemInfo) {
	memVir := MemInfo{}
	memInfoVir, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	memVir.Total = memInfoVir.Total
	memVir.Free = memInfoVir.Free
	memVir.Used = memInfoVir.Used
	memVir.UsedPercent = memInfoVir.UsedPercent
	//memInfoList = append(memInfoList, memVir)
	//memInfoSwap, err := mem.SwapMemory()
	//if err != nil {
	//	return
	//}
	//memVir.Total = memInfoSwap.Total
	//memVir.Free = memInfoSwap.Free
	//memVir.Used = memInfoSwap.Used
	//memVir.UsedPercent = memInfoSwap.UsedPercent
	//memInfoList = append(memInfoList, memVir)
	return memVir
}

// GetHostName 主机信息
func GetHostName() string {
	hostInfo, _ := host.Info()
	return hostInfo.Hostname
}

// GetIP
func GetIP() string {
	return middleware.LocalIp
}

// 磁盘信息
func GetDiskInfo() (diskInfoList []DiskInfo) {
	disks, err := disk.Partitions(true)
	if err != nil {
		return
	}
	for _, v := range disks {
		diskInfo := DiskInfo{}
		info, err := disk.Usage(v.Device)
		if err != nil {
			continue
		}
		diskInfo.Total = info.Total
		diskInfo.Free = info.Free
		diskInfo.Used = info.Used
		diskInfo.UsedPercent = info.UsedPercent
		diskInfoList = append(diskInfoList, diskInfo)
	}
	return
}

// GetNetwork 网络信息
func GetNetwork() (networks []Network) {
	netIOs, _ := net.IOCounters(true)
	if netIOs == nil {
		return
	}
	for _, netIO := range netIOs {
		var network = Network{}
		network.Name = netIO.Name
		network.BytesSent = netIO.BytesSent
		network.BytesRecv = netIO.BytesRecv
		network.PacketsSent = netIO.PacketsSent
		network.PacketsRecv = netIO.PacketsRecv
		networks = append(networks, network)
	}
	return
}

func InitLocalIp() {

	conn, err := gonet.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Logger.Error(fmt.Sprintf("udp服务：%s", err.Error()))
		return
	}
	localAddr := conn.LocalAddr().(*gonet.UDPAddr)
	middleware.LocalIp = strings.Split(localAddr.String(), ":")[0]
	log.Logger.Info("本机ip：", middleware.LocalIp)
}

func SendHeartBeatRedis() {
	timer := time.NewTicker(time.Duration(config.Conf.Heartbeat.Duration) * time.Second)
	for {
		select {
		case <-timer.C:
			CheckHeartBeat()
			hb, _ := json.Marshal(heartbeat)
			err := model.InsertHeartbeat(Key, middleware.LocalIp, string(hb))
			if err != nil {
				log.Logger.Error(fmt.Sprintf("机器ip:%s, 心跳发送失败, 写入redis失败:   %s", middleware.LocalIp, err.Error()))
			}
		}
	}
}

func SendMachineResources() {
	timer := time.NewTicker(time.Duration(config.Conf.Heartbeat.Resources) * time.Second)
	key := fmt.Sprintf("MachineMonitor:%s", middleware.LocalIp)
	for {
		select {
		case <-timer.C:
			hb, _ := json.Marshal(heartbeat)
			err := model.InsertMachineResources(key, string(hb))
			if err != nil {
				log.Logger.Error(fmt.Sprintf("机器ip:%s, 资源写入失败, 写入redis失败:   %s", middleware.LocalIp, err.Error()))
			}
		}
	}
}

func Logout() {
	field := middleware.LocalIp
	_ = model.DelMachine(Key, field)
	_ = model.DelResources(fmt.Sprintf("MachineMonitor:%s", middleware.LocalIp))
	log.Logger.Info("机器ip:", field, "注销成功")
}
