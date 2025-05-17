package logic

import (
	"encoding/json"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"scene/internal/biz/log"
	"scene/internal/conf"
	"scene/rao"
)

func GetEngineMap(ctx *gin.Context) (map[string]rao.HeartBeat, error) {
	req, err := http.NewRequest("GET", conf.Conf.Url.Machine+"/machine/available/list", nil)
	if err != nil {
		return nil, err
	}
	// 发起 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cp rao.CommonResponse[map[string]rao.HeartBeat]

	err = json.NewDecoder(resp.Body).Decode(&cp)
	if err != nil {
		log.Logger.Error("logic.engine.GetAvailableEngineMap.Decode，err:", err)
		return nil, err
	}
	return cp.Data, nil
}

func GetAvailableEngineList(ctx *gin.Context) ([]string, error) {
	// 获取所有施压机
	engineMap, err := GetEngineMap(ctx)
	if err != nil {
		log.Logger.Error("logic.engine.GetAvailableEngineList.Decode，err:", err)
		return nil, err
	}
	arr := make([]string, 0)
	// 获取使用中的施压机ip
	usingList, err := getUsingEngine(ctx)
	if err != nil {
		log.Logger.Error("logic.engine.GetAvailableEngineList.getUsingEngine，err:", err)
		return nil, err
	}
	for _, v := range engineMap {
		if !slice.Contain(usingList, v.IP) {
			arr = append(arr, v.IP)
		}
	}
	return arr, nil
}
