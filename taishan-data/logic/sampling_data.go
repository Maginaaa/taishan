package logic

import (
	"data/internal/biz/log"
	"data/internal/dal"
	"data/rao"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

func ClearSampling(ctx *gin.Context) {

	db := dal.MongoDB()
	// 列出所有集合名
	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Logger.Error("logic.sampling_data.ClearCache.ListCollectionNames，err:", err)
		return
	}

	// 打印集合名
	for _, collectionName := range collections {
		if !strings.Contains(collectionName, ":") {
			continue
		}
		begin := time.Now()
		log.Logger.Info("开始处理集合:", collectionName)
		arr := strings.Split(collectionName, ":")
		reportId := arr[0]
		sceneId, _ := strconv.Atoi(arr[1])
		cursor, _ := db.Collection(collectionName).Find(ctx, bson.D{{}}, options.Find().SetLimit(10000))
		var results []rao.HttpResponse
		err = cursor.All(ctx, &results)
		for _, res := range results {
			res.SceneId = int32(sceneId)
			db.Collection(reportId).InsertOne(ctx, res)
		}
		db.Collection(collectionName).Drop(ctx)
		log.Logger.Infof("集合: %s 处理完成, 耗时：%fs", collectionName, time.Now().Sub(begin).Seconds())

	}
}
